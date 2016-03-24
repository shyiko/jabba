package command

import (
	"io/ioutil"
	"encoding/json"
	"runtime"
	"net/http"
	"github.com/shyiko/jabba/cfg"
	"errors"
	"strconv"
)

type byOS map[string]byArch
type byArch map[string]byDistribution
type byDistribution map[string]map[string]string

func LsRemote() (map[string]string, error) {
	cnt, err := fetch(cfg.Index())
	if err != nil {
		return nil, err
	}
	var index byOS
	// todo: handle deserialization error
	json.Unmarshal(cnt, &index)
	return index[runtime.GOOS][runtime.GOARCH]["jdk"], nil
}

func fetch(url string) (content []byte, err error) {
	client := http.Client{Transport: RedirectTracer{}}
	res, err := client.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return nil, errors.New("GET " + url + " returned " + strconv.Itoa(res.StatusCode))
	}
	content, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	return
}
