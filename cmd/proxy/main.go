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
	"os"

	"github.com/google/tcpproxy"
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

	p := tcpproxy.Proxy{}
	p.AddRoute(proxyAddr, tcpproxy.To(serverAddr))
	log.Fatal(p.Run())
}
