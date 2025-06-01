FROM golang:alpine AS base
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
