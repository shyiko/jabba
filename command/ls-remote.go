package command

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Jabba-Team/jabba/cfg"
	"github.com/Jabba-Team/jabba/semver"
)

type byOS map[string]byArch
type byArch map[string]byDistribution
type byDistribution map[string]map[string]string

func LsRemote(os, arch string) (map[*semver.Version]string, error) {
	cnt, err := fetch(cfg.Index())
	if err != nil {
		return nil, err
	}
	var index byOS
	err = json.Unmarshal(cnt, &index)
	if err != nil {
		return nil, err
	}
	releaseMap := make(map[*semver.Version]string)
	for key, value := range index[os][arch] {
		var prefix string
		if key != "jdk" {
			if !strings.Contains(key, "@") {
				continue
			}
			prefix = key[strings.Index(key, "@")+1:] + "@"
		}
		for ver, url := range value {
			v, err := semver.ParseVersion(prefix + ver)
			if err != nil {
				return nil, err
			}
			releaseMap[v] = url
		}
	}
	return releaseMap, nil
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
