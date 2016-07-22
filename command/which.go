package command

import (
	"path/filepath"
	"github.com/shyiko/jabba/cfg"
)

func Which(selector string) (string, error) {
	aliasValue := GetAlias(selector)
	if aliasValue != "" {
		selector = aliasValue
	}
	ver, err := LsBestMatch(selector)
	if err != nil {
		return "", err
	}
	return filepath.Join(cfg.Dir(), "jdk", ver), nil
}
