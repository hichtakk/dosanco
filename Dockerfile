FROM ubuntu:18.04

COPY ./build/dosanco-apiserver_* /usr/local/bin/dosanco-apiserver
COPY ./build/dosanco_* /usr/local/bin/dosanco
COPY ./config.toml /etc/dosanco/config.toml

CMD ["/usr/local/bin/dosanco-apiserver"]