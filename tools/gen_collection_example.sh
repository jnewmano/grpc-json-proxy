#!/bin/sh

POSTMAN_COLLECTION_NAME="Hello World gRPC JSON API" protoc \
--doc_out=./ \
--doc_opt=./postman.tmpl,helloworld_collection.json \
helloworld.proto
