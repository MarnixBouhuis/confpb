package config

import (
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/ghodss/yaml"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ProtoMessage[T any] interface {
	*T
	proto.Message
}

type osFS struct{}

var _ fs.ReadFileFS = &osFS{}

func (*osFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

func (*osFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func FromFile[T any, TMsg ProtoMessage[T]](filename string) (TMsg, error) {
	return FromFileFS[T, TMsg](&osFS{}, filename)
}

func FromFileFS[T any, TMsg ProtoMessage[T]](fsys fs.FS, filename string) (TMsg, error) {
	bytes, err := fs.ReadFile(fsys, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filename, err)
	}

	switch path.Ext(filename) {
	case ".yaml", ".yml":
		return FromYAML[T, TMsg](bytes)
	case ".json":
		return FromJSON[T, TMsg](bytes)
	case ".pb":
		return FromPb[T, TMsg](bytes)
	case ".pb_text":
		return FromPbText[T, TMsg](bytes)
	default:
		return nil, fmt.Errorf("unknown file type: %s", filename)
	}
}

func FromYAML[T any, TMsg ProtoMessage[T]](bytes []byte) (TMsg, error) {
	// There is no direct proto unmarshaler for yaml files, convert the YAML to JSON and unmarshal that.
	json, err := yaml.YAMLToJSON(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert YAML to JSON: %w", err)
	}
	return FromJSON[T, TMsg](json)
}

func FromJSON[T any, TMsg ProtoMessage[T]](bytes []byte) (TMsg, error) {
	opts := &protojson.UnmarshalOptions{
		DiscardUnknown: false, // Error on unknown fields
	}

	a := &anypb.Any{}
	if err := opts.Unmarshal(bytes, a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON content to config: %w", err)
	}

	return anyToConfig[T, TMsg](a)
}

func FromPb[T any, TMsg ProtoMessage[T]](bytes []byte) (TMsg, error) {
	opts := &proto.UnmarshalOptions{
		DiscardUnknown: false, // Error on unknown fields
	}

	a := &anypb.Any{}
	if err := opts.Unmarshal(bytes, a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf binary content to config: %w", err)
	}

	return anyToConfig[T, TMsg](a)
}

func FromPbText[T any, TMsg ProtoMessage[T]](bytes []byte) (TMsg, error) {
	opts := &prototext.UnmarshalOptions{
		DiscardUnknown: false, // Error on unknown fields
	}

	a := &anypb.Any{}
	if err := opts.Unmarshal(bytes, a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf text content to config: %w", err)
	}

	return anyToConfig[T, TMsg](a)
}

func anyToConfig[T any, TMsg ProtoMessage[T]](a *anypb.Any) (TMsg, error) {
	opts := proto.UnmarshalOptions{
		DiscardUnknown: false, // Error on unknown fields
	}

	var conf T

	// Make sure *T is the same type as TMsg, otherwise anypb.UnmarshalTo does not accept the value.
	confPtr, ok := interface{}(&conf).(TMsg)
	if !ok {
		return nil, fmt.Errorf("type error: expected type %T to be a pointer to %T", confPtr, new(TMsg))
	}

	if err := anypb.UnmarshalTo(a, confPtr, opts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return confPtr, nil
}
