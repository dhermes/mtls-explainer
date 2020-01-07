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

FROM golang:1.13.5-alpine AS go_static

RUN apk --no-cache add git
RUN go get -u -v github.com/google/tcpproxy

COPY cmd/proxy/main.go /var/code/main.go
WORKDIR /var/code
RUN go build -o /usr/local/bin/proxy .

FROM alpine:3.10.3 AS app

COPY --from=go_static /usr/local/bin/proxy /usr/local/bin/proxy
ENTRYPOINT ["/usr/local/bin/proxy"]
