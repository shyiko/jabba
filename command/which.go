package command

import (
	"path"
	"github.com/shyiko/jabba/cfg"
)

func Which(ver string) (string, error) {
	resolved, err := resolveLocal(ver)
	if err != nil {
		return "", err
	}
	return path.Join(cfg.Dir(), "jdk", resolved), nil
}
