# confpb - Protobuf as config schema for Go
[![Go Reference](https://pkg.go.dev/badge/github.com/marnixbouhuis/confpb.svg)](https://pkg.go.dev/github.com/marnixbouhuis/confpb)
[![CI/CD Pipeline](https://github.com/MarnixBouhuis/confpb/actions/workflows/cicd.yaml/badge.svg)](https://github.com/MarnixBouhuis/confpb/actions/workflows/cicd.yaml)

**`confpb`** is a Go library that leverages Protocol Buffers (protobuf) to define, manage, and load configuration structures. It provides tools to map environment variables, set default values, merge configurations, and parse multiple formats like YAML, JSON, and protobuf files.

## Key Features

- `protoc-gen-env`: Automatically map environment variables to protobuf fields in Go code.
- `protoc-gen-default`: Define default values for protobuf fields directly in your `.proto` files.
- erging & Overlaying: Combine multiple protobuf files, allowing one to overlay values from another.
- Multi-format Parsing: Seamlessly parse configurations from various formats such as YAML, JSON, binary protobuf, and text protobuf.

## Why Use `confpb`?

- **Simplified configuration management:** Centralize and standardize configuration using protobuf, version your schemas, generate typed code from schemas.
- **Standardisation of config:** Use the wide ecosystem of protobuf tooling to lint your schemas, generate docs, detect breaking changes and validate values.
- **Flexibility:** environment binding, default generation and file loading are separated, allowing you to choose what features to implement.

---

## Quickstart Guide

### Example: Defining Your Configuration

Below is an example of how to define your application's configuration schema using protobuf. It includes environment variable mapping and default value definitions, though you can use either feature independently.

```protobuf
syntax = "proto3";
package config.v1;

import "confpb/v1/field.proto";  // Import confpb field options.

option go_package = "sample-app/gen/config/v1;configv1";

message ApplicationConfig {
  message ServerConfig {
     string host = 1 [(confpb.v1.env) = "HOST"];
     uint32 port = 2 [(confpb.v1.env) = "PORT"];
  }

  // Map environment variables:
  // SERVER_HOST=127.0.0.1
  // SERVER_PORT=8080
  ServerConfig server = 1 [(confpb.v1.env) = "SERVER"];

  string some_string = 2 [(confpb.v1.default).string = "foo"];

  // Set environment variables for repeated fields:
  // SOME_LIST_1 = "item1"
  // SOME_LIST_2 = "item2"
  // SOME_LIST_3 = "item3"
  // ...
  repeated string some_list = 3 [
    (confpb.v1.env) = "SOME_LIST",
    (confpb.v1.default).repeated_string = {
      values: ["default1", "default2"]
    }
  ];
}
```

### Loading Configuration in Go
Here’s how to load configuration in your Go application from default values, environment variables, and config files.

```go
package main

import (
    "github.com/marnixbouhuis/confpb/pkg/config"
    configv1 "sample-app/gen/config/v1"
)

func main() {
    // Start with default configuration
    conf := configv1.ApplicationConfigFromDefault()

    // Overlay with environment variables
    envConfig, err := configv1.ApplicationConfigFromEnv()
    if err != nil {
        panic(err)
    }
    conf = config.Overlay(conf, envConfig)

    // Overlay with configuration from a YAML file
    fileConfig, err := config.FromFile[configv1.ApplicationConfig]("./config.yaml")
    if err != nil {
        panic(err)
    }
    conf = config.Overlay(conf, fileConfig)

    // The config is now complete, with precedence in this order:
    // 1. Config file
    // 2. Environment variables
    // 3. Default values

    // Use the final configuration
    // fmt.Println(conf.SomeList[0])
}
```

Load values from a YAML file:
```
"@type": "type.googleapis.com/config.v1.ApplicationConfig"
server:
  host: 127.0.0.1
  port: 8443
some_string: example
some_list:
  - item1
  - item2
```

For a complete example, check out the [sample-app](https://github.com/MarnixBouhuis/confpb/tree/main/examples/sample-app).

---

## Installing the protoc plugins

There are two ways to install the generation plugins.
1. Download prebuilt binaries (Recommended)
2. Build from Source

Installing a prebuilt binary from release is recommended as this binary contains some extra version information.

### Option 1: Download prebuilt binaries
1. Head to the [releases page](https://github.com/MarnixBouhuis/confpb/releases) to download the appropriate binary for your platform.
   - Download both `protoc-gen-env` and `protoc-gen-default`.
2. Extract the downloaded archives.
3. Move the binaries to a directory in your `$PATH`, such as `/usr/local/bin` or `/bin`.

### Option 2: Building / installing from source
```bash
$ go install github.com/marnixbouhuis/confpb/cmd/protoc-gen-default@latest
$ go install github.com/marnixbouhuis/confpb/cmd/protoc-gen-env@latest
 ```
