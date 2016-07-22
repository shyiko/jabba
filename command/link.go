package command

import (
	"errors"
	"github.com/shyiko/jabba/cfg"
	"path/filepath"
	"os"
	"strings"
)

func Link(name string, dir string) error {
	if !strings.HasPrefix(name, "system@") {
		return errors.New("Name must begin with 'system@' (e.g. 'system@1.8.73')")
	}
	if dir == "" {
		return os.Remove(filepath.Join(cfg.Dir(), "jdk", name))
	} else {
		if err := assertJavaDistribution(dir); err != nil {
			return err
		}
		return os.Symlink(dir, filepath.Join(cfg.Dir(), "jdk", name))
	}
}

func GetLink(name string) string {
	res, err := os.Readlink(filepath.Join(cfg.Dir(), "jdk", name))
	if err != nil {
		return ""
	}
	return res
}
