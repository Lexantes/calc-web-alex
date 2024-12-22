FROM golang:1.23.4-alpine AS build

WORKDIR /app
COPY . .

ARG CGO_ENABLED=0 \
    GOARCH=amd64 \
    GOOS=linux \
    LDFLAGS_COMMON="-s -w"

RUN go build -buildmode=pie -ldflags "${LDFLAGS_COMMON}" -o ./bin/my-web-server ./cmd/main.go

FROM alpine:3.21.0 AS final

WORKDIR /my-app
COPY --from=build /app/bin/my-web-server ./my-web-server

ENTRYPOINT [ "/my-app/my-web-server" ]
CMD [ "--port", "8080" ]