FROM golang:alpine AS base

ARG APP_VERSION="development"

WORKDIR /build

COPY *go* /build
COPY internal /build/internal

RUN go build -ldflags="-X 'main.Version=${APP_VERSION}'" -o app

RUN apk add --no-cache --update ca-certificates tzdata

RUN adduser -D -H -h / -s /sbin/nologin app

FROM scratch AS packed

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /usr/share/zoneinfo/ /usr/share/zoneinfo/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /build/app /usr/local/bin/app

USER app

ENV TZ="Europe/Budapest"
EXPOSE 5353

ENTRYPOINT ["app"]