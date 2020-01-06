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

[1]: https://venilnoronha.io/a-step-by-step-guide-to-mtls-in-go
[2]: https://github.com/FiloSottile/mkcert
