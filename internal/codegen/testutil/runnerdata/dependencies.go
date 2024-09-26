package main

// Reference all possible dependencies that are used inside the tests so "go mod tidy" does not delete them from
// the go.mod file.

import (
	_ "github.com/stretchr/testify/assert"
	_ "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/structpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
)
