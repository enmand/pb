package main

import (
	"fmt"
	"strings"
)

// module represents a module that containing *.proto files
type module struct {
	Name        string
	Path        string
	CleanupFunc func() error
}

// Cleanup cleans up any temporary directories that were created
func (m *module) Cleanup() error {
	if m.CleanupFunc != nil {
		err := m.CleanupFunc()
		if err != nil {
			return err
		}
	}

	return nil
}

// Cleanup cleans up modules in this module set
func (m *modules) Cleanup() error {
	for _, mod := range *m {
		err := mod.Cleanup()
		if err != nil {
			return err
		}
	}

	return nil
}

// String returns the string representation of this module
func (m *module) String() string {
	return m.Path
}

type modules []module

// String returns the string representation of this module set. The moduleset
// is separated by `:`
func (m *modules) String() string {
	paths := []string{}
	for _, mod := range *m {
		paths = append(paths, mod.String())
	}

	return fmt.Sprintf("%s", strings.Join(paths, ":"))
}
