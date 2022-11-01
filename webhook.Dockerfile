FROM golang:1.19.2-alpine AS builder

RUN apk add --update --no-cache build-base git

RUN mkdir -p /build
WORKDIR /build

COPY . /build
RUN go mod download
RUN CGO_ENABLED=0 go install ./cmd/dia-webhook


FROM scratch

COPY --from=builder /go/bin/dia-webhook /usr/local/bin/dia-webhook

USER 10001

ENTRYPOINT ["/usr/local/bin/dia-webhook"]
