package testutil

import (
	"io"
	"io/fs"
	"strings"
	"testing"

	"github.com/jhump/protoreflect/desc/protoparse" //nolint:staticcheck // Disable deprecation check, we can't use buf protocompile directly because of dynamic extension types.
	"github.com/marnixbouhuis/confpb/internal/codegen"
	protofiles "github.com/marnixbouhuis/confpb/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gengo "google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type GenerationResult struct {
	Req  *pluginpb.CodeGeneratorRequest
	Resp *pluginpb.CodeGeneratorResponse
}

// RunGeneratorForFiles runs a codegen.FileGeneratorFunc for a specific set of proto files.
// Internally it creates a code generation requests and creates a new plugin instance.
// It then invokes the generator and returns the code generation response from the plugin.
func RunGeneratorForFiles(t *testing.T, generator codegen.FileGeneratorFunc, fs fs.FS, files ...string) *GenerationResult {
	t.Helper()

	// Generate codegen request, this is normally given by us by protoc via STDIN
	req := createCodegenRequestForFiles(t, fs, files)

	resp := runProtocPlugin(t, req, func(plugin *protogen.Plugin) {
		err := codegen.InvokeGeneratorForFiles(plugin, generator)
		require.NoError(t, err)
	})

	return &GenerationResult{
		Req:  req,
		Resp: resp,
	}
}

func createCodegenRequestForFiles(t *testing.T, filesFS fs.FS, filesToGenerate []string) *pluginpb.CodeGeneratorRequest {
	t.Helper()

	parser := protoparse.Parser{
		IncludeSourceCodeInfo: true,
		Accessor: func(filename string) (io.ReadCloser, error) {
			if strings.HasPrefix(filename, "confpb/") {
				return protofiles.Files.Open(filename)
			}
			return filesFS.Open(filename)
		},
	}

	// Build list of files to generate, order is important. Dependencies should go first, then files that depend on them.
	files := make([]string, 0, len(defaultLibraries)+len(filesToGenerate))
	files = append(files, defaultLibraries...)
	files = append(files, filesToGenerate...)

	// Parse proto files
	descriptors, err := parser.ParseFiles(files...)
	require.NoError(t, err)

	descriptorsProto := make([]*descriptorpb.FileDescriptorProto, 0, len(descriptors))
	for _, d := range descriptors {
		descriptorsProto = append(descriptorsProto, d.AsFileDescriptorProto())
	}

	pluginParams := ""
	return &pluginpb.CodeGeneratorRequest{
		FileToGenerate:        filesToGenerate,
		Parameter:             &pluginParams,
		ProtoFile:             descriptorsProto,
		SourceFileDescriptors: descriptorsProto,
		CompilerVersion:       &pluginpb.Version{},
	}
}

// generateFilesWithProtocGenGo generates go protobuf files using the same generation code as protoc-gen-go uses.
// It returns the generation response.
func generateFilesWithProtocGenGo(t *testing.T, req *pluginpb.CodeGeneratorRequest) *pluginpb.CodeGeneratorResponse {
	t.Helper()

	return runProtocPlugin(t, req, func(plugin *protogen.Plugin) {
		for _, file := range plugin.Files {
			if file.Generate {
				gengo.GenerateFile(plugin, file)
			}
		}
		plugin.SupportedFeatures = gengo.SupportedFeatures
		plugin.SupportedEditionsMaximum = gengo.SupportedEditionsMaximum
		plugin.SupportedEditionsMinimum = gengo.SupportedEditionsMinimum
	})
}

// runProtocPlugin creates a new protogen plugin instance based on a code generation request.
// It calls the callback function "fn" with this generator and returns the code generation response.
func runProtocPlugin(t *testing.T, req *pluginpb.CodeGeneratorRequest, fn func(plugin *protogen.Plugin)) *pluginpb.CodeGeneratorResponse {
	t.Helper()

	opts := &protogen.Options{}
	plugin, err := opts.New(req)
	require.NoError(t, err)

	fn(plugin)

	resp := plugin.Response()
	assert.Nil(t, resp.Error)

	return resp
}
