FROM alpine:lastest

RUN apk add --no-cache bash \
  make \
  curl \
  git

ENTRYPOINT ["/entrypoint.sh"]
CMD [ "-h" ]

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
