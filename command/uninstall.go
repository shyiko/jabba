package command

import (
	"os"
	"path/filepath"

	"github.com/Jabba-Team/jabba/cfg"
)

func Uninstall(selector string) error {
	ver, err := LsBestMatch(selector)
	if err != nil {
		return err
	}
	return os.RemoveAll(filepath.Join(cfg.Dir(), "jdk", ver))
}
