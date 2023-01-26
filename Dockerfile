FROM alpine:lastest

RUN apk add --no-cache bash \
  make \
  curl \
  git

