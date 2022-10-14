FROM gcr.io/test-servers-256610/ronin_ronin-base-image:base-image-0c47779be as builder
ENV GOPRIVATE="github.com/axieinfinity"
WORKDIR /opt

RUN --mount=type=secret,id=github_token git config --global url."https://x-access-token:$(cat /run/secrets/github_token)@github.com".insteadOf "https://github.com"

COPY . /opt/cleaner

RUN cd cleaner && make build
RUN cd cleaner && mkdir build && cp -p cleaner ./build

FROM debian:buster

RUN apt-get update  && apt-get install -y --no-install-recommends ca-certificates
RUN update-ca-certificates

WORKDIR /opt

COPY --from=builder /opt/cleaner/build /opt
CMD ["./cleaner"]
