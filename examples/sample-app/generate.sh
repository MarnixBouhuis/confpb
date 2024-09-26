#!/usr/bin/env bash

rm -rf ./gen/*

# Make sure to have the following things installed:
# - protoc (note some package managers ship really old version so protoc, make sure to have a recent one): https://grpc.io/docs/protoc-installation/
# - proto-gen-go: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# - protoc-gen-default: download / install from this repo's release page.
# - protoc-gen-env: download / install from this repo's release page.
protoc \
  --go_out=./gen/ \
  --go_opt=paths=source_relative \
  --default_out=./gen/ \
  --default_opt=paths=source_relative \
  --env_out=./gen/ \
  --env_opt=paths=source_relative \
  --proto_path=$PWD/../../proto/ \
  --proto_path=./proto/ \
  $(find ./proto/ -name "*.proto")
