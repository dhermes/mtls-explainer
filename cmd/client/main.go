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
// $ go run cmd/client/main.go \
// >   ./docker/tls-certs/rootCA-cert.pem  \
// >   ./docker/tls-certs/localhost-client-cert.pem \
// >   ./docker/tls-certs/localhost-client-key.pem

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

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

func getHostname() string {
	hostname, ok := os.LookupEnv("MTLS_SERVER_HOSTNAME")
	if ok {
		return hostname
	}

	return "localhost"
}

func main() {
	rootCAPath, certPath, keyPath, err := tlsKeyPaths(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// Read the key pair to create certificate
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create a CA certificate pool and add certificate to it
	caCert, err := ioutil.ReadFile(rootCAPath)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a HTTPS client and supply the created CA pool and certificate
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	// Request /hello via the created HTTPS client over port 8443 via GET
	hostname := getHostname()
	url := fmt.Sprintf("https://%s:8443/hello", hostname)
	r, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	// Read the response body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response body to stdout
	fmt.Printf("Received: %q\n", body)
}
