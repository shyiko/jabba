package cfg

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
)

func Dir() string {
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return path.Join(dir, ".jabba")
}

func Index() string {
	registry := os.Getenv("JABBA_INDEX")
	if registry == "" {
		registry = "https://github.com/shyiko/jabba/raw/master/index.json"
	}
	return registry
}
