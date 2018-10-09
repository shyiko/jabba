package command

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/shyiko/jabba/cfg"
)

func Exec(selector string, command string, commandArgs []string) error {
	aliasValue := GetAlias(selector)
	if aliasValue != "" {
		selector = aliasValue
	}
	ver, err := LsBestMatch(selector)
	if err != nil {
		return err
	}

	vars, err := GetVars(filepath.Join(cfg.Dir(), "jdk", ver))
	if err != nil {
		return err
	}

	for variableName, variableValue := range vars {
		os.Setenv(variableName, variableValue)
	}

	binary, err := exec.LookPath(command)
	if err != nil {
		return err
	}

	err = syscall.Exec(binary, commandArgs, os.Environ())
	return err
}
