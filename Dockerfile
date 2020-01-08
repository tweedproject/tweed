FROM golang:1.13 AS stage1

RUN apt-get update && \
    apt-get install libgpgme-dev libassuan-dev libbtrfs-dev libdevmapper-dev -y
WORKDIR /go/src/github.com/tweedproject/tweed
COPY . .

ENV GO111MODULE=on
RUN go build -mod=vendor ./cmd/tweed

FROM ubuntu:18.04 AS stage2
RUN apt-get update \
 && apt-get install -y curl unzip \
 && mkdir /bins \
 && curl -Lo /bins/safe       https://github.com/starkandwayne/safe/releases/download/v1.4.1/safe-linux-amd64 \
 && curl -Lo /bins/spruce     https://github.com/geofffranks/spruce/releases/download/v1.23.0/spruce-linux-amd64 \
 && curl -Lo /bins/jq         https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 \
 && chmod 0755 /bins/*

RUN curl -Lo /bins/duffle     https://github.com/cnabio/duffle/releases/download/0.3.5-beta.1/duffle-linux-amd64 \
    && chmod 0755 /bins/*

FROM ubuntu:18.04
RUN apt-get update \
 && apt-get install --no-install-recommends -y \
            libgpgme-dev libassuan-dev \
            ca-certificates -y # Needed by duffle to pull from docker hub \
            && rm -rf /var/lib/apt/lists/*
              # uuid-runtime curl net-tools \
              # postgresql \
              # redis \
              # mysql-client \
              # rabbitmq-server \
              # mongodb-clients \
              # kafkacat \

COPY --from=stage2 /bins/* /usr/bin/
COPY --from=stage1 /go/src/github.com/tweedproject/tweed/tweed /usr/bin
COPY bin      /tweed/bin

ADD entrypoint.sh /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD []
