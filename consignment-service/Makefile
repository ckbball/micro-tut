
build:
	protoc -I=$(GOPATH)/src/github.com/ckbball/micro-tut/consignment-service/proto/consignment --go_out=plugins=micro:$(GOPATH)/src/github.com/ckbball/micro-tut/consignment-service/proto/consignment/ $(GOPATH)/src/github.com/ckbball/micro-tut/consignment-service/proto/consignment/consignment.proto
	docker build -t consignment-service .

run:
	docker run -p 50051:50051 consignment-service