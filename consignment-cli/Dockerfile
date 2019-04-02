FROM golang:1.10.0 as builder

WORKDIR /go/src/github.com/ckbball/micro-tut/consignment-cli

COPY . .

RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -o consignment-cli -a -installsuffix cgo cli.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY consignment.json /app/consignment.json
COPY --from=builder /go/src/github.com/ckbball/micro-tut/consignment-cli/consignment-cli .

ENTRYPOINT ["./consignment-cli", "consignment.json"]
CMD ["./consignment-cli"]
