#!/bin/sh
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e -x

# Make sure CAROOT is set (it is used by `mkcert`)
if [[ -z "${CAROOT}" ]]; then
  echo "CAROOT environment variable should be set by the caller."
  exit 1
fi

# NOTE: `git` is needed for `go get`
apk --update --no-cache add git
go get -u -v github.com/FiloSottile/mkcert

# Clear out any existing root certificate (i.e. we want to always re-generate).
rm -f "${CAROOT}/rootCA-cert.pem" "${CAROOT}/rootCA-key.pem"

# (Re-)generate keys for `localhost`, store them **in** the CA root directory,
# which is expected to be a shared volume with the host.
cd "${CAROOT}"
mkcert \
  --client \
  --cert-file localhost-client-cert.pem \
  --key-file  localhost-client-key.pem \
  localhost
# NOTE: `server` and `proxy` are also valid hostnames so that the certificate
#       may be used within Dockerized networks where `localhost` can't be used
#       across containers.
mkcert \
  --client \
  --cert-file localhost-server-cert.pem \
  --key-file  localhost-server-key.pem \
  localhost server proxy

# Rename the root CA cert (for the benefit of the shared volume on the host).
mv "${CAROOT}/rootCA.pem" "${CAROOT}/rootCA-cert.pem"
