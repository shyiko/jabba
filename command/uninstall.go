package command

import (
	"github.com/shyiko/jabba/cfg"
	"path"
	"os"
)

func Uninstall(selector string) error {
	ver, err := LsBestMatch(selector)
	if err != nil {
		return err
	}
	return os.RemoveAll(path.Join(cfg.Dir(), "jdk", ver))
}
