#!/usr/bin/env sh

API_VERSION=0.1.3

rm -rf protobuf
mkdir protobuf
cd protobuf

# Download protoc
wget -q https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip
unzip protoc-3.5.1-linux-x86_64.zip
rm protoc-3.5.1-linux-x86_64.zip

# Download Veidemann API
wget -O - -q https://github.com/nlnwa/veidemann-api/archive/${API_VERSION}.tar.gz | tar --strip-components=2 -zx

go get -u github.com/golang/protobuf/proto
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go
