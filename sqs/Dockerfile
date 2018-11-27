FROM alpine:3.8

RUN apk --no-cache --update add ca-certificates

ADD bin/sqssource sqssource

RUN chmod +x sqssource

ENTRYPOINT ["./sqssource"]