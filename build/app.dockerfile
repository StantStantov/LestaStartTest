FROM golang:alpine AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
go build -v -o ./tmp/app ./cmd/...

FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /build/tmp/app ./tmp/app
COPY --from=builder /build/web/css/. ./web/css/.
ENTRYPOINT ["/app/tmp/app"]
