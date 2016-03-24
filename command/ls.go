package command

import (
	"io/ioutil"
	"github.com/shyiko/jabba/cfg"
	"github.com/shyiko/jabba/semver"
	"path"
)

var readDir = ioutil.ReadDir

// returns installed versions in descending order
func Ls() ([]string, error) {
	files, _ := readDir(path.Join(cfg.Dir(), "jdk"))
	var r []string
	for _, f := range files {
		if f.IsDir() {
			r = append(r, f.Name())
		}
	}
	return semver.Sort(r), nil
}
