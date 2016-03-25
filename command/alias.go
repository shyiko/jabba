package command

import (
	"errors"
	"io/ioutil"
	"github.com/shyiko/jabba/cfg"
	"path"
	"os"
)

func SetAlias(name string, ver string) (err error) {
	if name != "default" {
		return errors.New("At this point only 'default' alias is allowed")
	}
	if ver == "" {
		err = os.Remove(path.Join(cfg.Dir(), name + ".alias"))
	} else {
		err = ioutil.WriteFile(path.Join(cfg.Dir(), name + ".alias"), []byte(ver), 0666)
	}
	return
}

func GetAlias(name string) string {
	b, err := ioutil.ReadFile(path.Join(cfg.Dir(), name + ".alias"))
	if err != nil {
		return ""
	}
	return string(b)
}
