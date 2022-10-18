package command

import (
	"path/filepath"
	"runtime"

	"github.com/Jabba-Team/jabba/cfg"
)

func Which(selector string, home bool) (string, error) {
	aliasValue := GetAlias(selector)
	if aliasValue != "" {
		selector = aliasValue
	}
	ver, err := LsBestMatch(selector)
	if err != nil {
		return "", err
	}
	path := filepath.Join(cfg.Dir(), "jdk", ver)
	if home && runtime.GOOS == "darwin" {
		path = filepath.Join(path, "Contents", "Home")
	}
	return path, nil
}
