package includes

import (
	"fmt"
	"strings"
)

// Module represents a Module that containing *.proto files
type Module struct {
	Name        string
	Path        string
	CleanupFunc func() error
}

// Cleanup cleans up any temporary directories that were created
func (m *Module) Cleanup() error {
	if m.CleanupFunc != nil {
		err := m.CleanupFunc()
		if err != nil {
			return err
		}
	}

	return nil
}

// Cleanup cleans up modules in this module set
func (m *Modules) Cleanup() error {
	for _, mod := range *m {
		err := mod.Cleanup()
		if err != nil {
			return err
		}
	}

	return nil
}

// String returns the string representation of this module
func (m *Module) String() string {
	return m.Path
}

type Modules []Module

// String returns the string representation of this module set. The moduleset
// is separated by `:`
func (m *Modules) String() string {
	paths := []string{}
	for _, mod := range *m {
		paths = append(paths, mod.String())
	}

	return fmt.Sprintf("%s", strings.Join(paths, ":"))
}
