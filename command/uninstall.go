package command

import (
	"github.com/shyiko/jabba/cfg"
	"os"
	"path/filepath"
)

func Uninstall(selector string) error {
	ver, err := LsBestMatch(selector)
	if err != nil {
		return err
	}
	return os.RemoveAll(filepath.Join(cfg.Dir(), "jdk", ver))
}
