package command

import (
	"path"
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
	return path.Join(cfg.Dir(), "jdk", ver), nil
}
