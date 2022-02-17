FROM --platform=linux/amd64 ubuntu:20.04 AS builder

ENV GOVERSION 1.17.7

ENV GOROOT /opt/go
ENV GOPATH /go/src

ENV LANG en_US.UTF-8
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y git curl wget tzdata make gcc && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

RUN cd /opt && wget -nv --show-progress --progress=bar:force:noscroll https://golang.org/dl/go${GOVERSION}.linux-amd64.tar.gz && \
    tar zxf go${GOVERSION}.linux-amd64.tar.gz && rm go${GOVERSION}.linux-amd64.tar.gz && \
    ln -s /opt/go/bin/go /usr/bin/ && \
    mkdir -p ${GOPATH} && \
    mkdir /.cache && \
    chmod 777 /.cache

WORKDIR ${GOPATH}
