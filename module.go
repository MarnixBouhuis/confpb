package confpb

import (
	"embed"
)

// Embed go.mod and go.sum files, we need these in the e2e test runner.
//
//go:embed go.mod go.sum
var moduleFileFS embed.FS

// Mark variable as used. It's used inside internal/codegen/testutil/runner.go using "go:linkname".
// We access the variable this way because we don't want to expose it in the public API surface of this package.
var _ = moduleFileFS
