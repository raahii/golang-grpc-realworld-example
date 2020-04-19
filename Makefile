.PHONY: proto

proto: gen-service-proto gen-gateway-proto gen-swagger-proto

gen-service-proto:
	protoc \
		-I=/usr/local/include \
		-I=./proto \
		-I=${GOPATH}/src \
		-I=${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.3/third_party/googleapis/ \
		--go_out=plugins=grpc:./proto \
		./proto/*.proto

gen-gateway-proto:
	protoc \
		-I=/usr/local/include \
		-I=./proto \
		-I=${GOPATH}/src \
		-I=${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.3/third_party/googleapis/ \
		--grpc-gateway_out=logtostderr=true:./proto \
		./proto/*.proto

gen-swagger-proto:
	protoc \
		-I=/usr/local/include \
		-I=./proto \
		-I=${GOPATH}/src \
		-I=${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.3/third_party/googleapis/ \
		--swagger_out=logtostderr=true:./doc \
		./proto/*.proto
