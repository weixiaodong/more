
all: proto build

proto:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
# 	go get -u -v github.com/golang/protobuf/protoc-gen-go
# 	for file in $$(git ls-files '*.proto'); do \
# 		echo "protoc  $$file..."; \
# 		protoc -I $$(dirname $$file) --go_out=plugins=grpc:$$(dirname $$file) $$file; \
# 	done
	protoc -I protos --go_out=plugins=grpc:protos protos/*.proto;

build:
	go build .