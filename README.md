# mTLS

Code from [A step by step guide to mTLS in Go][1]

Certificates generated via the very helpful [FiloSottile/mkcert][2]
project, though done so in a container to avoid installing a fake root CA
onto the host. (This is insecure and should be done so with extreme caution.)

To generate a root CA and two distinct certificate pairs (for two different
services running on `localhost`):

```
./_bin/generate-tls-certs.sh
```

## Running on Host

In one shell:

```
go run cmd/server/main.go \
  ./docker/tls-certs/rootCA-cert.pem  \
  ./docker/tls-certs/localhost-server-cert.pem \
  ./docker/tls-certs/localhost-server-key.pem
```

and in another

```
go run cmd/client/main.go \
  ./docker/tls-certs/rootCA-cert.pem  \
  ./docker/tls-certs/localhost-client-cert.pem \
  ./docker/tls-certs/localhost-client-key.pem
```

## Running via `docker-compose`

To run both the client and server at once:

```
docker-compose \
  --file docker/docker-compose.yml \
  up \
  --build \
  --force-recreate \
  --abort-on-container-exit \
  --exit-code-from=server \
  --remove-orphans
```

This should produce output similar to:

```
...
Successfully tagged dhermes/mtls-explainer-proxy:latest
Creating docker_server_1 ... done
Creating docker_proxy_1  ... done
Creating docker_client_1 ... done
Attaching to docker_server_1, docker_proxy_1, docker_client_1
proxy_1   | Running TCP pass through proxy from :28443 to server:18443
server_1  | Listening on :18443
server_1  | Handling request for /hello
client_1  | Received: "Hello, world!\n"
docker_server_1 exited with code 0
Aborting on container exit...
Stopping docker_proxy_1  ... done
```

To clean up:

```
docker-compose \
  --file docker/docker-compose.yml \
  down
```

[1]: https://venilnoronha.io/a-step-by-step-guide-to-mtls-in-go
[2]: https://github.com/FiloSottile/mkcert



```
$ openssl x509 -in ./docker/tls-certs/localhost-client-cert.pem -serial -noout
serial=CC6880577A039A6185EDFD28D5698429
$ openssl x509 -in ./docker/tls-certs/localhost-server-cert.pem -serial -noout
serial=67ADFF08C978DA5FDB11058AE70E1534
```


Used:

https://redflagsecurity.net/2019/03/10/decrypting-tls-wireshark/ <-- uses `sslkeylog.log` to get session keys to decrypt
Also added `localhost-server-key.pem` as an RSA key to the SSL protocol with `http` as the interceptor

https://hpbn.co/transport-layer-security-tls/
https://www.cloudflare.com/learning/ssl/what-happens-in-a-tls-handshake/
https://blog.cloudflare.com/logjam-the-latest-tls-vulnerability-explained/
https://www.thesslstore.com/blog/explaining-ssl-handshake/
