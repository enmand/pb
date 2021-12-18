package includes

// Local loads local paths with the name in the key and the local path in the
// value
func Local(path map[string]string) (Modules, error) {
	ms := Modules{}
	for name, path := range path {
		ms = append(ms, Module{Name: name, Path: path})
	}

	return ms, nil
}
