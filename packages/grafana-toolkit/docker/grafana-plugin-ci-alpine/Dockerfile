FROM alpine:3.15

USER root

COPY scripts scripts
COPY install /usr/local

WORKDIR scripts

RUN ./deploy.sh
