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

set -e

# ``readlink -f`` is not our friend on macOS.
SCRIPT_FILE=$(python -c "import os; print(os.path.realpath('${0}'))")
BIN_DIR=$(dirname ${SCRIPT_FILE})
ROOT_DIR=$(dirname ${BIN_DIR})

docker run \
  --rm \
  --volume "${ROOT_DIR}/_bin/generate-tls-certs-on-alpine.sh":/bin/generate-tls-certs-on-alpine.sh \
  --volume "${ROOT_DIR}/docker/tls-certs:/var/tls-certs" \
  --env CAROOT=/var/tls-certs \
  golang:1.13.5-alpine \
  /bin/generate-tls-certs-on-alpine.sh
