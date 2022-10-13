#syntax=docker/dockerfile:experimental
FROM gcr.io/test-servers-256610/ronin_ronin-base-image:base-image-0c47779be as builder
WORKDIR /opt

ENV GO111MODULE=on
ENV GOPRIVATE="github.com/axieinfinity"
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"
RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

COPY . /opt/cleaner
RUN --mount=type=ssh,mode=741,uid=100,gid=102 cd cleaner && make build
RUN cd cleaner && mkdir build && cp -p cleaner ./build

FROM debian:buster

RUN apt-get update  && apt-get install -y --no-install-recommends ca-certificates
RUN update-ca-certificates

WORKDIR /opt

COPY --from=builder /opt/cleaner/build /opt
CMD ["./cleaner"]
