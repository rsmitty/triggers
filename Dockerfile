FROM alpine:3.8

RUN apk --no-cache --update add ca-certificates

ADD bin/gcssource gcssource

RUN chmod +x gcssource

ENTRYPOINT ["./gcssource"]