.PHONY: all build clean proto test run

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
BINARY_NAME=aliasme
BINARY_UNIX=$(BINARY_NAME)_unix

# Proto files
PROTO_FILES=proto/service.proto

all: proto build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v .

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# proto:
# 	protoc -I . \
# 		--go_out . --go_opt paths=source_relative \
# 		--go-grpc_out . --go-grpc_opt paths=source_relative \
# 		--grpc-gateway_out . --grpc-gateway_opt paths=source_relative \
# 		$(PROTO_FILES)

test:
	$(GOTEST) -v ./...

run: build
	./$(BINARY_NAME)

deps:
	$(GOGET) -v ./...
	$(GOMOD) tidy

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/server

PROTOC_IMAGE := registry.hub.docker.com/rvolosatovs/protoc:4

#-----------------------------------------------------------------------------
# BUILD
#-----------------------------------------------------------------------------

.PHONY: proto
proto:
	podman run --rm -v $$(pwd):$$(pwd) -w $$(pwd) $(PROTOC_IMAGE) \
		-I$(GOPATH)/src \
		-Iproto \
		--go_out=$(GOPATH)/src \
		--go-grpc_out=$(GOPATH)/src \
		--grpc-gateway_out=logtostderr=true:$(GOPATH)/src \
		--openapiv2_out=allow_merge=true,merge_file_name=gostatus,output_format=yaml,allow_repeated_fields_in_body=true:pkg/static/swagger/ \
		--openapiv2_out=allow_merge=true,merge_file_name=gostatus,output_format=json,allow_repeated_fields_in_body=true:pkg/static/swagger/ \
		--validate_out="lang=go:$(GOPATH)/src" \
		./proto/*.proto
