package command

import (
	"io/ioutil"
	"github.com/shyiko/jabba/cfg"
	"github.com/shyiko/jabba/semver"
	"path"
	"sort"
	"fmt"
	"os"
)

var readDir = ioutil.ReadDir

func Ls() ([]*semver.Version, error) {
	files, _ := readDir(path.Join(cfg.Dir(), "jdk"))
	var r []*semver.Version
	for _, f := range files {
		if f.IsDir() || f.Mode() & os.ModeSymlink == os.ModeSymlink {
			v, err := semver.ParseVersion(f.Name())
			if err != nil {
				return nil, err
			}
			r = append(r, v)
		}
	}
	sort.Sort(sort.Reverse(semver.VersionSlice(r)))
	return r, nil
}

func LsBestMatch(selector string) (ver string, err error) {
	local, err := Ls()
	if err != nil {
		return
	}
	rng, err := semver.ParseRange(selector)
	if err != nil {
		return
	}
	for _, v := range local {
		if rng.Contains(v) {
			ver = v.String()
			break
		}
	}
	if ver == "" {
		err = fmt.Errorf("%s isn't installed", rng)
	}
	return
}
