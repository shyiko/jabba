package command

import (
	"github.com/shyiko/jabba/cfg"
	"os"
	"path/filepath"
	"regexp"
)

func Deactivate() ([]string, error) {
	pth, _ := os.LookupEnv("PATH")
	plSep := string(os.PathListSeparator)
	rgxp := regexp.MustCompile(regexp.QuoteMeta(filepath.Join(cfg.Dir(), "jdk")) + "[^" + plSep + "]+[" + plSep + "]")
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
