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
Creating network "docker_service-net" with the default driver
Creating docker_server_1 ... done
Creating docker_client_1 ... done
Attaching to docker_server_1, docker_client_1
server_1  | Listening on :8443
server_1  | Handling request for /hello
client_1  | Received: "Hello, world!\n"
docker_client_1 exited with code 0
Aborting on container exit...
Stopping docker_server_1 ... done
```

To clean up:

```
docker-compose \
  --file docker/docker-compose.yml \
  down
```

[1]: https://venilnoronha.io/a-step-by-step-guide-to-mtls-in-go
[2]: https://github.com/FiloSottile/mkcert
