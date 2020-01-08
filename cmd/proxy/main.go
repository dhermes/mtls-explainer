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

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/google/tcpproxy"
)

// NOTE: Ensure that
//       * Listener satisfies net.Listener
//       * Conn satisfies net.Conn
var (
	_ net.Conn     = (*Conn)(nil)
	_ net.Listener = (*Listener)(nil)
)

// envOrDefault returns an environment variable at `name` if provided, or falls
// back to `defaultPort`.
func envOrDefault(name, defaultPort string) string {
	port, ok := os.LookupEnv(name)
	if ok {
		return port
	}

	return defaultPort
}

// Conn embeds a `net.Conn` and "spies" on reads and writes.
type Conn struct {
	original net.Conn
}

// Read forwards `Read()` and "spies" on the content that was read.
func (c *Conn) Read(b []byte) (int, error) {
	n, err := c.original.Read(b)
	if err == nil {
		fmt.Printf("Read() completed %d bytes: %x\n", n, b[:n])
		// fmt.Printf("Read() completed %d bytes\n", n)
	} else {
		fmt.Printf("Read() failed with %v\n", err)
	}
	return n, err
}

// Write forwards `Write()` and "spies" on the content that was written.
func (c *Conn) Write(b []byte) (int, error) {
	n, err := c.original.Write(b)
	if err == nil {
		fmt.Printf("Write() completed %d bytes: %x\n", n, b[:n])
		// fmt.Printf("Write() completed %d bytes\n", n)
	} else {
		fmt.Printf("Write() failed with %v\n", err)
	}
	return n, err
}

// Close forwards `Close()`.
func (c *Conn) Close() error {
	return c.original.Close()
}

// LocalAddr forwards `LocalAddr()`.
func (c *Conn) LocalAddr() net.Addr {
	return c.original.LocalAddr()
}

// RemoteAddr forwards `RemoteAddr()`.
func (c *Conn) RemoteAddr() net.Addr {
	return c.original.RemoteAddr()
}

// SetDeadline forwards `SetDeadline()`.
func (c *Conn) SetDeadline(t time.Time) error {
	return c.original.SetDeadline(t)
}

// SetReadDeadline forwards `SetReadDeadline()`.
func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.original.SetReadDeadline(t)
}

// SetWriteDeadline forwards `SetWriteDeadline()`.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.original.SetWriteDeadline(t)
}

// Listener embeds a net.Listener, it forwards `Close()` and `Addr()` and uses
// `Accept()` to return a wrapped `net.Conn`.
type Listener struct {
	original net.Listener
}

// Accept forwards `Accept()` and wraps the returned connection.
func (ln *Listener) Accept() (net.Conn, error) {
	conn, err := ln.original.Accept()
	fmt.Printf("conn type: %T\n", conn)
	wrapped := &Conn{original: conn}
	return wrapped, err
}

// Close forwards `Close()`.
func (ln *Listener) Close() error {
	return ln.original.Close()
}

// Addr forwards `Addr()`.
func (ln *Listener) Addr() net.Addr {
	return ln.original.Addr()
}

func listenFunc(network, address string) (net.Listener, error) {
	if network != "tcp" {
		return nil, fmt.Errorf("Unexpected network %q", network)
	}

	ln, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	fmt.Printf("ln type: %T\n", ln)
	wrapped := &Listener{original: ln}
	return wrapped, nil
}

func main() {
	proxyAddr := fmt.Sprintf(
		"%s:%s",
		envOrDefault("PROXY_SERVER_HOSTNAME", ""),
		envOrDefault("PROXY_SERVER_PORT", "9090"),
	)
	serverAddr := fmt.Sprintf(
		"%s:%s",
		envOrDefault("PROXIED_HOSTNAME", "localhost"),
		envOrDefault("PROXIED_PORT", "8080"),
	)

	fmt.Printf(
		"Running TCP pass through proxy from %s to %s\n",
		proxyAddr,
		serverAddr,
	)

	p := tcpproxy.Proxy{ListenFunc: listenFunc}
	p.AddRoute(proxyAddr, tcpproxy.To(serverAddr))
	log.Fatal(p.Run())
}

// http.Server
// Handler -> defaults to http.DefaultServeMux
// TLSConfig
// ConnContext <-- to mess with ctx on a new connection (ServerContextKey)
// ConnState <-- fun to mess with
