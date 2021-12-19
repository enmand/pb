package includes

import "github.com/unerror/id-hub/tools/protoc/config"

// Local loads local paths with the name in the key and the local path in the
// value
func Local(dep config.Dependency) (Module, error) {
	return Module{Name: dep.Name, Path: dep.Version}, nil
}
