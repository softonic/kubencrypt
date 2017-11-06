#------------------------------------------------------------------------------
# Set the base image for subsequent instructions:
#------------------------------------------------------------------------------

FROM alpine:3.6

#------------------------------------------------------------------------------
# Build-time arguments:
#------------------------------------------------------------------------------

ARG SDK_VERSION="178.0.0"

#------------------------------------------------------------------------------
# Environment variables:
#------------------------------------------------------------------------------

ENV GOPATH="/go"
ENV PATH="/opt/google-cloud-sdk/bin:${PATH}"

#------------------------------------------------------------------------------
# Build and install:
#------------------------------------------------------------------------------

RUN SDK_URL="https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${SDK_VERSION}-linux-x86_64.tar.gz"; \
    apk add -U --no-cache -t dev git go musl-dev curl \
    && apk add -U --no-cache python \
    && mkdir /opt && curl -L ${SDK_URL} | tar zx -C /opt \
    && gcloud config set component_manager/disable_update_check true \
    && gcloud config set core/disable_usage_reporting true \
    && gcloud config set metrics/environment github_docker_image \
    && go get github.com/softonic/kubencrypt/cmd/kubencrypt \
    && cp ${GOPATH}/bin/kubencrypt /usr/local/bin \
    && apk del --purge dev && rm -rf /tmp/* /go

#------------------------------------------------------------------------------
# Entrypoint:
#------------------------------------------------------------------------------

ENTRYPOINT [ "kubencrypt" ]
