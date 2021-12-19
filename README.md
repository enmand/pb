# pb

`pb` is a Protocol Buffers Build tool that manages dependencies and build configuration
for `protoc`.

# Flags

- `proto-path` - Path to `.proto` to compile using configuration
- `config` - Path to configuration file to use for compilation

# Config

Configuration is provided in HCL with the following top-level configuration fields:

## `go`

`go` configures `protoc` dependencies based on the go.mod, and includes any modules
that contain `*.proto` files. Fields are:

- `path` - Path to the `go.mod` file. Relative to configuration or absolute
- `ignores` - array of directory paths to ignore for importing `*.proto` files

### Example

```hcl
go {
  path = "go.mod"
  ignores = ["internal"]
}
```

## `dependency(type, label)`

`dependencies` configures the `label` dependencies for the provided `type` at the
`version` provided

The fields are:

- `type` (label) - The dependency type (either: `git` or `local`)
- `name` (label) - The name and location of the dependency
- `version` - The tag or commit to use as a version

### Example

```hcl
dependency "git" "github.com/googleapis/googleapis" {
    version = "5518740a67d22a7ad1b0b7656c98211ca5e19307"
}
```

## `plugin(name)`

`plugin` configures `protoc` plugins for code generation. The fields are:

- `name` - The binary name suffix (i.e. protoc-gen-\<name\>)
- `path` - The path to place generated files, relative to the configuration
- `rel_path` - The path to place generated files, relative to the `*.proto` file being generated for
- `options` - Options to provide the `protoc` plugin

### Example

```hcl
plugin "go" {
    rel_path = "."
    options = {
        "paths": "source_relative",
    }
}

plugin "go-grpc" {
    rel_path = "."
    options = {
        "paths": "source_relative",
    }
}
```

# Full example

```hcl
go {
	path = "./go.mod"
	ignore = ["vendor", "test", "example", "internal"]
}

dependency "git" "github.com/googleapis/googleapis" {
	version = "5518740a67d22a7ad1b0b7656c98211ca5e19307"
}

plugin "go" {
	rel_path = "."
	options = {
		"paths": "source_relative",
	}
}

plugin "go-grpc" {
	rel_path = "."
	options = {
		"paths": "source_relative",
	}
}

plugin "swagger" {
	path = "./genspec/openapi/"
	options = {
		"allow_merge": "true",
		"merge_file_name": "api"
	}
}

plugin "grpc-gateway" {
	rel_path = "."
	options = {
		"paths": "source_relative",
		"allow_patch_feature": "false",
	}
}
```

# Using with `go generate`

You can use this tool with `go generate` in place of using Makefiles for protoc generate, by adding a build-tag to a file alongside the `.proto` file:

```go
//go:generate go run github.com/enmand/pb -proto-path=path/to/service.proto --config=path/to/proto.hcl
package service
```
