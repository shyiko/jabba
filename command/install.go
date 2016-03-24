package command

import (
	"os/exec"
	"fmt"
	"runtime"
	"errors"
	"sort"
	"strings"
	"os"
	"io/ioutil"
	"path"
	"net/http"
	"io"
	"github.com/shyiko/jabba/cfg"
	"github.com/shyiko/jabba/semver"
	wmark "github.com/wmark/semver"
	log "github.com/Sirupsen/logrus"
	"regexp"
	"github.com/mitchellh/ioprogress"
)

func Install(qualifier string) (ver string, err error) {
	var releaseMap map[string]string
	if strings.Contains(qualifier, "=") {
		// <version>=<url>
		split := strings.SplitN(qualifier, "=", 2)
		// <version> has to be valid per semver
		_, err = wmark.NewVersion(split[0])
		if err != nil {
			return
		}
		qualifier = split[0]
		ver = qualifier
		releaseMap = map[string]string{qualifier: split[1]}
	}
	// check whether it's already installed
	local, err := Ls()
	if err != nil {
		return
	}
	i := sort.Search(len(local), func(i int) bool {
		return local[i] <= qualifier
	})
	if i < len(local) && local[i] == qualifier {
		// already installed
		return qualifier, nil
	}
	if releaseMap == nil {
		rng, err := wmark.NewRange(qualifier)
		if err != nil {
			return "", err
		}
		releaseMap, err = LsRemote()
		if err != nil {
			return "", err
		}
		var vs = make([]string, len(releaseMap))
		var i = 0
		for k := range releaseMap {
			vs[i] = k
			i++
		}
		vs = semver.Sort(vs)
		for i := range vs {
			v, _ := wmark.NewVersion(vs[i])
			if rng.Contains(v) {
				ver = vs[i]
				break
			}
		}
		if ver == "" {
			return ver, errors.New("No compatible version found for " + qualifier +
			"\nValid install targets: " + strings.Join(vs, ", "))
		}
	}
	url := releaseMap[ver]
	if matched, _ := regexp.MatchString("^\\w+[+]\\w+://", url); !matched {
		return ver, errors.New("URL must contain qualifier, e.g. tgz+http://...")
	}
	var fileType string = url[0:strings.Index(url, "+")]
	url = url[strings.Index(url, "+") + 1:len(url)]
	var file string
	var deleteFileWhenFinnished bool
	if strings.HasPrefix(url, "file://") {
		file = strings.TrimPrefix(url, "file://")
	} else {
		log.Info("Downloading ", ver, " (", url, ")")
		file, err = download(url)
		if err != nil {
			return
		}
		deleteFileWhenFinnished = true
	}
	switch runtime.GOOS {
	case "darwin":
		err = installOnDarwin(ver, file, fileType)
	case "linux":
		err = installOnLinux(ver, file, fileType)
	default:
		err = errors.New(runtime.GOOS + " OS is not supported")
	}
	if err == nil && deleteFileWhenFinnished {
		os.Remove(file)
	}
	return
}

type RedirectTracer struct {
	Transport http.RoundTripper
}

func (self RedirectTracer) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	transport := self.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	resp, err = transport.RoundTrip(req)
	if err != nil {
		return
	}
	switch resp.StatusCode {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect:
		log.Debug("Following ", resp.StatusCode, " redirect to ", resp.Header.Get("Location"))
	}
	return
}

func download(url string) (file string, err error) {
	tmp, err := ioutil.TempFile("", "jabba-d-")
	if err != nil {
		return
	}
	file = tmp.Name()
	log.Debug("Saving ", url, " to ", file)
	// todo: timeout
	client := http.Client{Transport: RedirectTracer{}}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return fmt.Errorf("too many redirects")
		}
		if len(via) != 0 {
			// https://github.com/golang/go/issues/4800
			for attr, val := range via[0].Header {
				if _, ok := req.Header[attr]; !ok {
					req.Header[attr] = val
				}
			}
		}
		return nil
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", "oraclelicense=accept-securebackup-cookie")
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	progressTracker := &ioprogress.Reader{
		Reader: res.Body,
		Size: res.ContentLength,
	}
	_, err = io.Copy(tmp, progressTracker)
	if err != nil {
		return
	}
	return
}

func installOnDarwin(ver string, file string, fileType string) error {
	if fileType != "dmg" {
		return errors.New(fileType + " is not supported")
	}
	tmp, err := ioutil.TempDir("", "jabba-i-")
	if err != nil {
		return err
	}
	basename := path.Base(file)
	mountpoint := tmp + "/" + basename
	target := cfg.Dir() + "/jdk/" + ver
	err = executeSH([][]string{
		[]string{"Mounting " + file, "hdiutil mount -mountpoint " + mountpoint + " " + file},
		[]string{"Extracting " + file + " to " + target,
			"pkgutil --expand " + mountpoint + "/*.pkg " + tmp + "/" + basename + "-pkg"},
		[]string{"", "mkdir -p " + target},

		// oracle
		[]string{"", "if [ -f " + tmp + "/" + basename + "-pkg/jdk*.pkg/Payload" + " ]; then " +
		"tar xvf " + tmp + "/" + basename + "-pkg/jdk*.pkg/Payload -C " + target +
		"; fi"},

		// apple
		[]string{"", "if [ -f " + tmp + "/" + basename + "-pkg/JavaForOSX.pkg/Payload" + " ]; then " +
		"tar -xzf " + tmp + "/" + basename + "-pkg/JavaForOSX.pkg/Payload -C " + tmp + "/" + basename + "-pkg &&" +
		"mv " + tmp + "/" + basename + "-pkg/Library/Java/JavaVirtualMachines/*/Contents " + target + "/Contents" +
		"; fi"},

		[]string{"Unmounting " + file, "hdiutil unmount " + mountpoint},
	})
	if err == nil {
		if _, err := os.Stat(target + "/Contents/Home/bin/java"); os.IsNotExist(err) {
			err = errors.New("Unsupported DMG structure. " +
			"Please open a ticket at https://github.com/shyiko/jabba/issue " +
			"(specify URI you tried to install)")
		}
	}
	if err != nil {
		// remove target ~/.jabba/jdk/<version>
		os.RemoveAll(target)
	} else {
		os.RemoveAll(tmp)
	}
	return err
}

func installOnLinux(ver string, file string, fileType string) (err error) {
	target := cfg.Dir() + "/jdk/" + ver
	var cmd [][]string
	var tmp string
	switch fileType {
	case "bin":
		tmp, err = ioutil.TempDir("", "jabba-i-")
		if err != nil {
			return
		}
		cmd = [][]string{
			[]string{"", "mv " + file + " " + tmp},
			[]string{"Extracting " + path.Join(tmp, path.Base(file)) + " to " + target,
				"cd " + tmp + " && echo | sh " + path.Base(file) + " && mv jdk*/ " + target},
		}
	case "tgz":
		cmd = [][]string{
			[]string{"", "mkdir -p " + target},
			[]string{"Extracting " + file + " to " + target,
				"tar xvf " + file + " --strip-components=1 -C " + target},
		}
	default:
		return errors.New(fileType + " is not supported")
	}
	err = executeSH(cmd)
	if err != nil {
		// remove target ~/.jabba/jdk/<version>
		os.RemoveAll(target)
	} else {
		if tmp != "" {
			os.RemoveAll(tmp)
		}
	}
	return
}

func executeSH(cmd [][]string) error {
	for _, command := range cmd {
		if command[0] != "" {
			log.Info(command[0])
		}
		out, err := exec.Command("sh", "-c", command[1]).CombinedOutput()
		if err != nil {
			log.Error(string(out))
			return errors.New("'" + command[1] + "' failed: " + err.Error())
		}
	}
	return nil
}
