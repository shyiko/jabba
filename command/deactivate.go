package command

import (
	"os"
	"github.com/shyiko/jabba/cfg"
	"path"
	"regexp"
)

func Deactivate() ([]string, error) {
	pth, _ := os.LookupEnv("PATH")
	rgxp := regexp.MustCompile(regexp.QuoteMeta(path.Join(cfg.Dir(), "jdk")) + "[^:]+[:]")
	// strip references to ~/.jabba/jdk/*, otherwise leave unchanged
	pth = rgxp.ReplaceAllString(pth, "")
	javaHome, overrideWasSet := os.LookupEnv("JAVA_HOME_BEFORE_JABBA")
	if !overrideWasSet {
		javaHome, _ = os.LookupEnv("JAVA_HOME")
	}
	return []string{
		"export PATH=" + pth,
		"export JAVA_HOME=" + javaHome,
		"unset JAVA_HOME_BEFORE_JABBA",
	}, nil
}
