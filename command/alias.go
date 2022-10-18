package command

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Jabba-Team/jabba/cfg"
)

func SetAlias(name string, ver string) (err error) {
	if ver == "" {
		err = os.Remove(filepath.Join(cfg.Dir(), name+".alias"))
	} else {
		err = ioutil.WriteFile(filepath.Join(cfg.Dir(), name+".alias"), []byte(ver), 0666)
	}
	return
}

func GetAlias(name string) string {
	b, err := ioutil.ReadFile(filepath.Join(cfg.Dir(), name+".alias"))
	if err != nil {
		return ""
	}
	return string(b)
}
