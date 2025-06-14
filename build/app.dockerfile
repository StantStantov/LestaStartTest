FROM golang:alpine AS base
RUN apk add git
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

FROM base AS test

FROM base AS build
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
  go build -v -o ./tmp/app ./cmd/...

FROM alpine:latest AS runner
WORKDIR /app
COPY --from=build /src/tmp/app ./tmp/app
ENTRYPOINT ["/app/tmp/app"]
