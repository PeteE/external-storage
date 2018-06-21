FROM alpine:3.6
RUN apk add --no-cache ca-certificates

COPY freenas-provisioner /usr/local/bin/freenas-provisioner

ENV USERNAME "dummy"
ENV PASSWORD "dummy"
ENV ADDRESS "192.168.1.10"
ENV POOL "wheatpool"
ENV SCHEME "https"
ENTRYPOINT ["/usr/local/bin/freenas-provisioner", "start"]
