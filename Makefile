.PHONY: proto

proto:
	protoc \
		-I=/usr/local/include \
		-I=./proto \
		-I=${GOPATH}/src \
		-I=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:./proto \
		--grpc-gateway_out=logtostderr=true:./proto \
		--swagger_out=logtostderr=true:./doc \
		./proto/*.proto

test:
	go test -v ./handler -parallel 4
