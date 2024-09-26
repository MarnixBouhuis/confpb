package testutil

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	// Needed for importing the moduleFileFS from the module.go in the root of this repo.
	_ "unsafe"

	// Needed for importing the moduleFileFS from the module.go in the root of this repo.
	_ "github.com/marnixbouhuis/confpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

// DefaultLibraries is a list of default proto files to include when creating a codegen request.
// The order of these paths is important, files that depend on other files should be defined after the original
// dependency. If A imports B and B imports C, then the order should be: C, B, A.
var defaultLibraries = []string{
	"google/protobuf/any.proto",
	"google/protobuf/source_context.proto",
	"google/protobuf/type.proto",
	"google/protobuf/api.proto",
	"google/protobuf/descriptor.proto",
	"google/protobuf/duration.proto",
	"google/protobuf/empty.proto",
	"google/protobuf/field_mask.proto",
	"google/protobuf/timestamp.proto",
	"google/protobuf/struct.proto",
	"google/protobuf/wrappers.proto",
	"confpb/v1/field.proto", // Special case, we later auto resolve paths starting with "confpb/" to the proto files in this repo.
}

//go:embed runnerdata
var runnerTemplateFS embed.FS

// Load go.mod / go.sum from module root. This variable is private (we don't want to expose this in the public API
// surface) so we have to use go:linkname.
//
//go:linkname moduleFileFS github.com/marnixbouhuis/confpb.moduleFileFS
var moduleFileFS embed.FS

// testEvent is a JSON event emitted by running `go test -json ./...`.
type testEvent struct {
	Time    time.Time `json:"Time"`    //nolint:tagliatelle // This is how test2json outputs data
	Action  string    `json:"Action"`  //nolint:tagliatelle // This is how test2json outputs data
	Package string    `json:"Package"` //nolint:tagliatelle // This is how test2json outputs data
	Test    string    `json:"Test"`    //nolint:tagliatelle // This is how test2json outputs data
	Elapsed float64   `json:"Elapsed"` //nolint:tagliatelle // This is how test2json outputs data
	Output  string    `json:"Output"`  //nolint:tagliatelle // This is how test2json outputs data
}

type e2eTestInfo struct {
	Failed bool
	Output string
}

func RunTestInE2ERunner(t *testing.T, res *GenerationResult, testCode string) {
	t.Helper()

	// Create a temporary directory for the test files to live in
	tempDir, err := os.MkdirTemp("", "confpb_e2e_test_")
	require.NoError(t, err)
	t.Logf("Created temporary directory for E2E test runner: %s", tempDir)

	// Make sure temporary directory is cleaned up after we're done running the test
	t.Cleanup(func() {
		// Allow skipping cleanup of failed test cases
		if os.Getenv("CONFPB_E2E_TEST_DIR_CLEANUP") == "false" && t.Failed() {
			t.Logf("Skipping cleanup of E2E test runner temproray directory since CONFPB_E2E_TEST_DIR_CLEANUP is false")
			return
		}
		require.NoError(t, os.RemoveAll(tempDir))
		t.Logf("Deleted temporary directory: %s", tempDir)
		t.Logf("You can disable cleanup of test case files by setting the environment variable \"CONFPB_E2E_TEST_DIR_CLEANUP\" to \"false\"")
	})

	// Copy go.mod and go.sum files from moduleFileFS into the temporary directory
	addModuleFilesToTestRunner(t, tempDir)

	// Copy all files from the runnerTemplateFS into the temporary directory
	runnerFiles, err := fs.Sub(runnerTemplateFS, "runnerdata")
	require.NoError(t, err)
	CopyAllFromFS(t, runnerFiles, tempDir)

	// Copy generated files into temporary directory
	for _, file := range res.Resp.File {
		destPath := filepath.Join(tempDir, *file.Name)
		err = os.WriteFile(destPath, []byte(*file.Content), 0o600)
		require.NoError(t, err)
		t.Logf("Wrote generated file to %s", destPath)
	}

	// Generate code files using "protoc-gen-go", our generator extends code generated by this generator.
	for _, file := range generateFilesWithProtocGenGo(t, res.Req).File {
		destPath := filepath.Join(tempDir, *file.Name)
		err = os.WriteFile(destPath, []byte(*file.Content), 0o600)
		require.NoError(t, err)
		t.Logf("Wrote generated protoc-gen-go file to %s", destPath)
	}

	// Create the `main.go` file containing the testCode inside the temporary directory
	err = os.WriteFile(filepath.Join(tempDir, "main_test.go"), []byte(testCode), 0o600)
	require.NoError(t, err)

	// Run the tests using `go test`
	cmd := exec.Command("go", "test", "-v", "-json", "./...")
	cmd.Dir = tempDir
	out, runnerErr := cmd.CombinedOutput()

	// First try to parse the test result
	var executedTests bool
	resultsPerTest := map[string]*e2eTestInfo{}
	decoder := json.NewDecoder(bytes.NewReader(out))
	for decoder.More() {
		var event testEvent
		require.NoError(t, decoder.Decode(&event))

		if event.Action == "run" {
			executedTests = true
		}

		testID := event.Package + "/" + event.Test
		testInfo, ok := resultsPerTest[testID]
		if !ok {
			testInfo = &e2eTestInfo{}
		}

		if event.Action == "fail" {
			testInfo.Failed = true
		}

		testInfo.Output += event.Output

		resultsPerTest[testID] = testInfo
	}

	// Make sure to error if we did not run any tests
	if !executedTests {
		t.Fatal("No tests where found / executed")
	}

	// Log each test output in its own t.Run group, this way we get grouping of tests running in the e2e runner.
	for testID, testInfo := range resultsPerTest {
		t.Run(testID, func(t *testing.T) {
			t.Helper()
			t.Log(testInfo.Output)
			if testInfo.Failed {
				t.Fail()
			}
		})
	}

	// Runner output parsed, now check if we exited with a non 0 status code
	require.NoError(t, runnerErr)
}

func addModuleFilesToTestRunner(t *testing.T, path string) {
	t.Helper()

	// Copy go.sum file from moduleFileFS to runner directory
	err := CopyFileFromFSToDisk(moduleFileFS, "go.sum", filepath.Join(path, "go.sum"))
	require.NoError(t, err)

	// Copy go.mod file from moduleFileFS to runner directory
	modFilePath := filepath.Join(path, "go.mod")
	err = CopyFileFromFSToDisk(moduleFileFS, "go.mod", modFilePath)
	require.NoError(t, err)

	// Update the go.mod file in the runner directory
	b, err := os.ReadFile(modFilePath)
	require.NoError(t, err)

	content := string(b)
	// Change the name of the module, so it does not conflict when importing things like the runtime package
	content = strings.ReplaceAll(content, "module github.com/marnixbouhuis/confpb", "module e2e_test_runner")
	// Add replace to the github.com/marnixbouhuis/confpb on localdisk
	content += "require github.com/marnixbouhuis/confpb v0.0.0-00010101000000-000000000000\n"
	content += fmt.Sprintf("replace github.com/marnixbouhuis/confpb => %s\n", getModulePath(t))

	err = os.WriteFile(modFilePath, []byte(content), 0o600)
	require.NoError(t, err)
}

func getModulePath(t *testing.T) string {
	t.Helper()

	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedModule}
	pkgs, err := packages.Load(cfg, ".")
	require.NoError(t, err, "Failed to load package info, are you running the tests from within the module itself?")

	if len(pkgs) != 1 || pkgs[0].Module == nil {
		t.Fatal("Failed to load package info (no package / module info), are you running the tests from within the module?")
	}

	// Make sure we are running in the right package
	assert.Equal(t, "github.com/marnixbouhuis/confpb", pkgs[0].Module.Path)
	return pkgs[0].Module.Dir
}
