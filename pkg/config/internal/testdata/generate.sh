#!/usr/bin/env bash
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  testdata.proto
