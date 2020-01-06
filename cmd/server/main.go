// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Thanks to Venil Noronha (https://github.com/venilnoronha), the original
// author of this code:
// https://venilnoronha.io/a-step-by-step-guide-to-mtls-in-go

// Expected to be run from the root of the repository via
// $ go run cmd/server/main.go \
// >   ./docker/tls-certs/rootCA-cert.pem  \
// >   ./docker/tls-certs/localhost-server-cert.pem \
// >   ./docker/tls-certs/localhost-server-key.pem

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func makeHelloHandler(beginShutdown chan struct{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello, world!\n")
		close(beginShutdown)
	}
}

func shutdownServer(server *http.Server, beginShutdown, shutdownComplete chan struct{}) {
	<-beginShutdown

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Failed shutting down: %v", err)
	}

	close(shutdownComplete)
}

func tlsKeyPaths(args []string) (rootCAPath, certPath, keyPath string, err error) {
	// TODO: Use real flags, this is too ad-hoc.
	if len(args) != 4 {
		err = fmt.Errorf("Expected exactly three arguments (CA cert filepath, cert filepath, key filepath)")
		return
	}

	relPath := args[1]
	rootCAPath, err = filepath.Abs(relPath)
	if err != nil {
		err = fmt.Errorf("Failed to get absolute path for %q: %w", relPath, err)
		return
	}

	relPath = args[2]
	certPath, err = filepath.Abs(relPath)
	if err != nil {
		err = fmt.Errorf("Failed to get absolute path for %q: %w", relPath, err)
		return
	}

	relPath = args[3]
	keyPath, err = filepath.Abs(relPath)
	if err != nil {
		err = fmt.Errorf("Failed to get absolute path for %q: %w", relPath, err)
		return
	}

	return
}

func main() {
	rootCAPath, certPath, keyPath, err := tlsKeyPaths(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// Channels to coordinate server shutdown.
	beginShutdown := make(chan struct{})
	shutdownComplete := make(chan struct{})

	// Set up a /hello resource handler
	http.HandleFunc("/hello", makeHelloHandler(beginShutdown))

	// Create a CA certificate pool and add certificate to it
	caCert, err := ioutil.ReadFile(rootCAPath)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}
	go shutdownServer(server, beginShutdown, shutdownComplete)

	// Listen to HTTPS connections with the server certificate and wait
	if err := server.ListenAndServeTLS(certPath, keyPath); err != http.ErrServerClosed {
		log.Fatalf("Failed listen and serve: %v", err)
	}

	<-shutdownComplete
}
