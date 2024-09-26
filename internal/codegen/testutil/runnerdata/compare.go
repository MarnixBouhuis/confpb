package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func protoEqual(t *testing.T, expected proto.Message, actual proto.Message) {
	t.Helper()
	if diff := cmp.Diff(expected, actual, protocmp.Transform()); diff != "" {
		t.Errorf("Different protobuf messages (- expected, + actual):\n%s", diff)
	}
}

// Mark protoEqual as used, it's used in generated test code.
var _ = protoEqual
