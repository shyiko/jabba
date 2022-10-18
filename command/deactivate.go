package command

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/Jabba-Team/jabba/cfg"
)

func Deactivate() ([]string, error) {
	pth, _ := os.LookupEnv("PATH")
	rgxp := regexp.MustCompile(regexp.QuoteMeta(filepath.Join(cfg.Dir(), "jdk")) + "[^:]+[:]")
	// strip references to ~/.jabba/jdk/*, otherwise leave unchanged
	pth = rgxp.ReplaceAllString(pth, "")
	javaHome, overrideWasSet := os.LookupEnv("JAVA_HOME_BEFORE_JABBA")
	if !overrideWasSet {
		javaHome, _ = os.LookupEnv("JAVA_HOME")
	}
	return []string{
		"export PATH=\"" + pth + "\"",
		"export JAVA_HOME=\"" + javaHome + "\"",
		"unset JAVA_HOME_BEFORE_JABBA",
	}, nil
}
