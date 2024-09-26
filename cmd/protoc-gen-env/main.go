package main

import (
	"github.com/marnixbouhuis/confpb/internal/codegen"
	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
)

func main() {
	codegen.RunProtocPlugin(envgen.GenerateFile)
}
