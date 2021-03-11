FROM golang:1.16-alpine AS build
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -v -o ./bin/resourceserver ./cmd/resourceserver/*.go

FROM alpine
WORKDIR /app
COPY --from=build /build/bin/resourceserver .
ENV APP_PUBLICKEYFILE="/cert/cert.pem" \
    APP_WS_BINDADDRESS=":80"
EXPOSE $APP_WS_BINDADDRESS
ENTRYPOINT ["/app/resourceserver"]