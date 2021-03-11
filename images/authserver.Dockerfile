FROM golang:1.16-alpine AS build
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -v -o ./bin/authserver ./cmd/authserver/*.go

FROM alpine
WORKDIR /app
COPY --from=build /build/bin/authserver .
ENV APP_PRIVATEKEYFILE="/cert/key.pem" \
    APP_WS_BINDADDRESS=":80"
EXPOSE $APP_WS_BINDADDRESS
ENTRYPOINT ["/app/authserver"]