# Test runner
This folder contains code used for running code generation e2e tests.

In some e2e tests we generate proto files, these files will then be copied into this folder together with a test case.
This test case will then be executed by the test runner.

The `go.mod` and `go.sum` files are copied from the root of this module into the e2e temp folder so the same dependencies
are available as in these files.
