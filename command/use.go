package command

import (
	"os"
	"github.com/shyiko/jabba/cfg"
	"path"
	"runtime"
	"regexp"
)

func Use(ver string) ([]string, error) {
	aliasValue := GetAlias(ver)
	if aliasValue != "" {
		ver = aliasValue
	}
	resolved, err := resolveLocal(ver)
	if err != nil {
		return nil, err
	}
	pth, _ := os.LookupEnv("PATH")
	rgxp := regexp.MustCompile(regexp.QuoteMeta(path.Join(cfg.Dir(), "jdk")) + "[^:]+[:]")
	// strip references to ~/.jabba/jdk/*, otherwise leave unchanged
	pth = rgxp.ReplaceAllString(pth, "")
	javaHome := path.Join(cfg.Dir(), "jdk", resolved)
	if runtime.GOOS == "darwin" {
		javaHome = path.Join(javaHome, "Contents", "Home")
	}
	systemJavaHome, overrideWasSet := os.LookupEnv("JAVA_HOME_BEFORE_JABBA")
	if !overrideWasSet {
		systemJavaHome, _ = os.LookupEnv("JAVA_HOME")
	}
	return []string{
		"export PATH=" + path.Join(javaHome, "bin") + ":" + pth,
		"export JAVA_HOME=" + javaHome,
		"export JAVA_HOME_BEFORE_JABBA=" + systemJavaHome,
	}, nil
}
