package runtime

import "github.com/marnixbouhuis/confpb/internal/version"

const (
	// Version is the current version of this runtime.
	// Generated code will use this value to determine if this runtime is compatible or not.
	Version = version.RuntimeVersion

	// MinimumSupportedCodegenVersion is the minimum supported codegen version by this runtime.
	// This value can be incremented to drop support for old generated runtime code.
	MinimumSupportedCodegenVersion = 1
)

// EnforceVersion is used as a compile time check in generated code to make sure the runtime version is compatible
// with the codegen version.
// This works the same way as how protoc-gen-go uses EnforceVersion.
// If this EnforceVersion causes compilation to fail it means that:
// - The runtime is too old, the generated code required a newer runtime. Update the runtime.
// - The generated code is too old for this runtime version. Generate the code with a new version of the generator.
type EnforceVersion uint

// Make sure minimum supported version is always lower or equal to runtime version.
const (
	_ = EnforceVersion(Version - MinimumSupportedCodegenVersion)
)
