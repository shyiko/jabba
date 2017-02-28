package command

import (
	"github.com/shyiko/jabba/cfg"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var lookPath = exec.LookPath

func Current() string {
	javaPath, err := lookPath("java")
	if err == nil {
		prefix := filepath.Join(cfg.Dir(), "jdk") + string(os.PathSeparator)
		if strings.HasPrefix(javaPath, prefix) {
			index := strings.Index(javaPath[len(prefix):], string(os.PathSeparator))
			if index != -1 {
				return javaPath[len(prefix) : len(prefix)+index]
			}
		}
	}
	return ""
}
