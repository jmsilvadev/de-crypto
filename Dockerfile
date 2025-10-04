FROM golang:1.24-alpine AS build
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /app/de-crypto ./cmd/main.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app/de-crypto /usr/local/bin/de-crypto
COPY data ./data
CMD ["de-crypto"]
