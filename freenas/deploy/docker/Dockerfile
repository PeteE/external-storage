FROM alpine:3.6
RUN apk add --no-cache ca-certificates

COPY freenas-provisioner /usr/local/bin/freenas-provisioner

ENTRYPOINT ["/usr/local/bin/freenas-provisioner", "start"]
