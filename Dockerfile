FROM golang:1.13 AS stage1

RUN apt-get update && \
    apt-get install libgpgme-dev libassuan-dev \
    libbtrfs-dev libdevmapper-dev -y
WORKDIR /go/src/github.com/tweedproject/tweed
COPY . .

ENV GO111MODULE=on
RUN go build -mod=vendor ./cmd/tweed

FROM ubuntu:18.04 AS stage2
RUN apt-get update && \
    apt-get install -y software-properties-common && \
    add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update \
    && apt-get install -y curl unzip \
    && mkdir /bins \
    && curl -Lo /bins/safe       https://github.com/starkandwayne/safe/releases/download/v1.4.1/safe-linux-amd64 \
    && curl -Lo /bins/spruce     https://github.com/geofffranks/spruce/releases/download/v1.23.0/spruce-linux-amd64 \
    && curl -Lo /bins/jq         https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 \
    && curl -Lo /bins/runc       https://github.com/opencontainers/runc/releases/download/v1.0.0-rc9/runc.amd64 \
    && chmod 0755 /bins/*

RUN apt-get install -y git make golang-go \
    libgpgme-dev libassuan-dev libdevmapper-dev

RUN export GOPATH=/tmp/gopath && \
    git clone https://github.com/containers/skopeo $GOPATH/src/github.com/containers/skopeo && \
    cd $GOPATH/src/github.com/containers/skopeo && make binary-local && \
    cp skopeo /bins/skopeo && rm -r $GOPATH

FROM ubuntu:18.04
RUN apt-get update \
    && apt-get install --no-install-recommends -y \
    libgpgme-dev libassuan-dev libdevmapper-dev \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=stage2 /bins/* /usr/bin/
COPY --from=stage1 /go/src/github.com/tweedproject/tweed/tweed /usr/bin
COPY bin      /tweed/bin

ADD entrypoint.sh /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD []
