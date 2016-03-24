package command

import (
	"os"
	"github.com/shyiko/jabba/cfg"
	"path"
	"runtime"
	"regexp"
)

func Use(ver string) ([]string, error) {
	resolved, err := resolveLocal(ver)
	if err != nil {
		return nil, err
	}
	rgxp := regexp.MustCompile(regexp.QuoteMeta(path.Join(cfg.Dir(), "jdk")) + "[^:]+[:]")
	p, _ := os.LookupEnv("PATH")
	p = rgxp.ReplaceAllString(p, "")
	javaHome := path.Join(cfg.Dir(), "jdk", resolved)
	if runtime.GOOS == "darwin" {
		javaHome = path.Join(javaHome, "Contents", "Home")
	}
	return []string{
		"export PATH=" + path.Join(javaHome, "bin") + ":" + p,
		"export JAVA_HOME=" + javaHome,
	}, nil
}
