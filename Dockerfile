FROM golang:1.20.4-alpine as build
RUN apk add make
RUN mkdir -p /app/src
WORKDIR /app/src

COPY go.mod . 
COPY go.sum . 
RUN go mod download
COPY . .
RUN make build

FROM alpine:3.7
COPY --from=build /app/src/main /root/main
COPY tinkoff_config.yaml /root/tinkoff_config.yaml

WORKDIR /root/

CMD ["./main"]