package main

import (
	"github.com/marnixbouhuis/confpb/internal/codegen"
	"github.com/marnixbouhuis/confpb/internal/codegen/defaultgen"
)

func main() {
	codegen.RunProtocPlugin(defaultgen.GenerateFile)
}
