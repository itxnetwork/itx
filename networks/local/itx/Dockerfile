FROM --platform=linux/arm64 golang:1.20-bullseye as build-env

# Install minimum necessary dependencies
ENV CGO_ENABLED=1
ARG ARCH=amd64

ENV PACKAGES curl make git libc-dev bash gcc-aarch64-linux-gnu
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y $PACKAGES

# Set working directory for the build
WORKDIR /go/src/github.com/itxnetwork/itx

# Add source files
COPY . .

# build Ethermint
RUN make build-linux arch=$ARCH

# Final image
FROM golang:1.20-bullseye as final

WORKDIR /

RUN apt-get update

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/itxnetwork/itx/build/itxd /
COPY --from=build-env /go/src/github.com/itxnetwork/itx/scripts/start-docker.sh /

EXPOSE 26656 26657 1317 9090 8545 8546

# Run itxd by default, omit entrypoint to ease using container with itxd
ENTRYPOINT ["/bin/bash", "-c"]