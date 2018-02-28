#!/usr/bin/env sh

rm -rf veidemann_api html/swagger/*
mkdir veidemann_api
cd protobuf

bin/protoc -I. \
  --go_out=plugins=grpc:../veidemann_api \
  --grpc-gateway_out=logtostderr=true:../veidemann_api \
  --swagger_out=logtostderr=true:../html/swagger \
  *.proto
