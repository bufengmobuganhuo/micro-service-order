

.PHONY: proto
proto:
	sudo docker run --rm -v $(shell pwd):$(shell pwd) -w $(shell pwd) -e ICODE=3FC927246EC68EAD cap1573/cap-protoc -I ./ --micro_out=./ --go_out=./ ./proto/order/order.proto

.PHONY: build
build: 

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o order-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t order-service:latest
