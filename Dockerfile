FROM golang:1.20.2-bullseye AS build-env

WORKDIR /go/src/github.com/itxnetwork/itx

RUN apt-get update -y
RUN apt-get install git -y

COPY . .

RUN make build

FROM golang:1.20.2-bullseye

RUN apt-get update -y
RUN apt-get install ca-certificates jq -y

WORKDIR /root

COPY --from=build-env /go/src/github.com/itxnetwork/itx/build/itxd /usr/bin/itxd

EXPOSE 26656 26657 1317 9090 8545 8546

CMD ["itxd"]
