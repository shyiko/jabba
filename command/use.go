package command

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/shyiko/jabba/cfg"
)

func Use(selector string) ([]string, error) {
	aliasValue := GetAlias(selector)
	if aliasValue != "" {
		selector = aliasValue
	}
	ver, err := LsBestMatch(selector)
	if err != nil {
		return nil, err
	}
	return usePath(filepath.Join(cfg.Dir(), "jdk", ver))
}

func GetVars(path string) (map[string]string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	pth, _ := os.LookupEnv("PATH")
	rgxp := regexp.MustCompile(regexp.QuoteMeta(filepath.Join(cfg.Dir(), "jdk")) + "[^:]+[:]")
	// strip references to ~/.jabba/jdk/*, otherwise leave unchanged
	pth = rgxp.ReplaceAllString(pth, "")
	if runtime.GOOS == "darwin" {
		path = filepath.Join(path, "Contents", "Home")
	}
	systemJavaHome, overrideWasSet := os.LookupEnv("JAVA_HOME_BEFORE_JABBA")
	if !overrideWasSet {
		systemJavaHome, _ = os.LookupEnv("JAVA_HOME")
	}

	paths := make(map[string]string)
	paths["PATH"] = filepath.Join(path, "bin") + string(os.PathListSeparator) + pth
	paths["JAVA_HOME"] = path
	paths["JAVA_HOME_BEFORE_JABBA"] = systemJavaHome

	return paths, nil
}

func usePath(path string) ([]string, error) {
	vars, err := GetVars(path)
	if err != nil {
		return nil, err
	}

	paths := []string{}

	for variableName, variableValue := range vars {
		paths = append(paths, "export "+variableName+"=\""+variableValue+"\"")
	}

	return paths, nil
}
