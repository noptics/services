.PHONY: proto

nats-proxy:
	docker build -t nats-proxy -f ./cmd/nats-proxy/Dockerfile .

registry:
	docker build -t noptics-registry -f ./cmd/registry/Dockerfile .

proto:
	protoc -I ./pkg/protos ./pkg/protos/*.proto --go_out=plugins=grpc:./pkg/nproto
