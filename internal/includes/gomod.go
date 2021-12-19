package includes

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jbowes/cling"
	"github.com/enmand/pb/internal/config"
)

// goModules returns the modules that contain *.proto files
func GoModDependencies(cfg config.Go) (Modules, error) {
	mods := exec.Command("go", "list", "-f", "{{.Path}}={{.Dir}}", "-m", "all")
	mods.Dir = filepath.Dir(cfg.Path)
	mods.Stderr = os.Stderr
	out, err := mods.Output()
	if err != nil {
		return nil, cling.Wrap(err, "unable to get go modules")
	}

	protoMods := Modules{}
	for _, mod := range strings.Split(string(out), "\n") {
		if mod == "" {
			continue
		}

		parts := strings.Split(mod, "=")
		if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
			continue
		}

		files := []string{}
		err := filepath.Walk(parts[1], func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if cfg.Ignores != nil {
				for _, ignore := range *cfg.Ignores {
					if strings.Contains(path, ignore) {
						return nil
					}
				}
			}

			if filepath.Ext(path) == ".proto" {
				files = append(files, path)
			}

			return nil
		})
		if err != nil {
			return nil, cling.Wrap(err, "unable to walk module")
		}

		if len(files) > 0 {
			mod, err := linkMod(parts[0], parts[1])
			if err != nil {
				return nil, cling.Wrap(err, "unable to link module")
			}

			protoMods = append(protoMods, mod...)
		}
	}

	return protoMods, nil
}

// linkTmp links the module into a temporary directory, and returns both the
// module with a relative path, and a module with a fully-resolved path
func linkMod(name, modPath string) (Modules, error) {
	tmp, err := ioutil.TempDir("", "protoc-")
	if err != nil {
		return nil, cling.Wrap(err, "unable to create temporary directory")
	}

	path := filepath.Join(tmp, name)
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return nil, cling.Wrap(err, "unable to create module dir")
	}

	err = os.Symlink(modPath, path)
	if err != nil {
		return nil, cling.Wrap(err, "unable to symlink module")
	}

	return Modules{
		{
			Name:        name,
			Path:        modPath,
			CleanupFunc: nil,
		},
		{
			Name: name,
			Path: tmp,
			CleanupFunc: func() error {
				return os.RemoveAll(tmp)
			},
		},
	}, nil
}
