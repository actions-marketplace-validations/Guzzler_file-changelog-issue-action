FROM golang:1.14 AS builder
COPY . /src/github.com/guzzler/hello-github-actions
WORKDIR /src/github.com/guzzler/hello-github-actions
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on \
  go build \
  -a \
  -o /bin/pr-changelog \
  /src/github.com/guzzler/hello-github-actions

FROM alpine:3.9 as certs-installer
RUN apk add --update ca-certificates

FROM scratch
COPY --from=builder /bin/pr-size /bin/pr-changelog
COPY --from=certs-installer /etc/ssl/certs /etc/ssl/certs
ENTRYPOINT ["/bin/pr-changelog"]
CMD [""]