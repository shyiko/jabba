package command

import (
	"errors"
	"github.com/shyiko/jabba/cfg"
	"github.com/shyiko/jabba/semver"
	"path/filepath"
	"os"
	"strings"
)

func Link(selector string, dir string) error {
	if !strings.HasPrefix(selector, "system@") {
		return errors.New("Name must begin with 'system@' (e.g. 'system@1.8.73')")
	}
	// <version> has to be valid per semver
	if _, err := semver.ParseVersion(selector); err != nil {
		return err
	}
	if dir == "" {
		ver, err := LsBestMatch(selector)
		if err != nil {
			return err
		}
		return os.Remove(filepath.Join(cfg.Dir(), "jdk", ver))
	} else {
		if err := assertJavaDistribution(dir); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join(cfg.Dir(), "jdk"), 0755); err != nil {
			return err
		}
		return os.Symlink(dir, filepath.Join(cfg.Dir(), "jdk", selector))
	}
}

func GetLink(name string) string {
	res, err := os.Readlink(filepath.Join(cfg.Dir(), "jdk", name))
	if err != nil {
		return ""
	}
	return res
}
