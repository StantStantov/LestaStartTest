FROM golang:alpine AS base
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

FROM base AS build
COPY . .
RUN --mount=target=. \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -v -o /tmp/app ./cmd/...

FROM alpine:latest AS runner
WORKDIR /app
COPY --from=build /src/web/css/. ./web/css/.
COPY --from=build /tmp/app ./tmp/app
ENTRYPOINT ["/app/tmp/app"]
