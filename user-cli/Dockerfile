FROM golang:1.10.0 as builder

WORKDIR /go/src/github.com/ckbball/micro-tut/user-cli

COPY . .

RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -o user-cli -a -installsuffix cgo cli.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/ckbball/micro-tut/user-cli/user-cli .

ENTRYPOINT ["./user-cli"]
CMD ["./user-cli"]