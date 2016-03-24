package command

import (
	"github.com/shyiko/jabba/cfg"
	"path"
	"os"
	"sort"
	"github.com/wmark/semver"
	"errors"
)

func Uninstall(ver string) error {
	resolved, err := resolveLocal(ver)
	if err != nil {
		return err
	}
	return os.RemoveAll(path.Join(cfg.Dir(), "jdk", resolved))
}

func resolveLocal(ver string) (string, error) {
	local, err := Ls()
	if err != nil {
		return "", err
	}
	i := sort.Search(len(local), func(i int) bool {
		return local[i] <= ver
	})
	var resolved string
	if i < len(local) && local[i] == ver {
		resolved = ver
	}
	if resolved == "" {
		// ver might be a range
		rng, err := semver.NewRange(ver)
		if err != nil {
			return "", err
		}
		for i := range local {
			v, _ := semver.NewVersion(local[i])
			if rng.Contains(v) {
				resolved = local[i]
				break
			}
		}
	}
	if resolved == "" {
		return "", errors.New(ver + " isn't installed")
	}
	return resolved, nil
}
