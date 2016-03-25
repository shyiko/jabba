package command

import (
	"os/exec"
	"strings"
	"github.com/shyiko/jabba/cfg"
	"path"
)

var lookPath = exec.LookPath

func Current() string {
	javaPath, err := lookPath("java")
	if err == nil {
		prefix := path.Join(cfg.Dir(), "jdk") + "/"
		if strings.HasPrefix(javaPath, prefix) {
			index := strings.Index(javaPath[len(prefix):], "/")
			if index != -1 {
				return javaPath[len(prefix):len(prefix) + index]
			}
		}
	}
	return ""
}
