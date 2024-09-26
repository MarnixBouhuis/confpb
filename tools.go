//go:build tools
// +build tools

package confpb

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser/v2"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
