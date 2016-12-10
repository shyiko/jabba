package command

import (
	"os/exec"
	"fmt"
	"runtime"
	"errors"
	"strings"
	"os"
	"io/ioutil"
	"path/filepath"
	"net/http"
	"io"
	"github.com/shyiko/jabba/cfg"
	"github.com/shyiko/jabba/semver"
	log "github.com/Sirupsen/logrus"
	"regexp"
	"github.com/mitchellh/ioprogress"
	"sort"
	"archive/zip"
	"github.com/shyiko/jabba/w32"
)

func Install(selector string) (string, error) {
	var releaseMap map[*semver.Version]string
	var ver *semver.Version
	var err error
	// selector can be in form of <version>=<url>
	if strings.Contains(selector, "=") && strings.Contains(selector, "://") {
		split := strings.SplitN(selector, "=", 2)
		selector = split[0]
		// <version> has to be valid per semver
		ver, err = semver.ParseVersion(selector)
		if err != nil {
			return "", err
		}
		releaseMap = map[*semver.Version]string{ver: split[1]}
	} else {
		// ... or a version (range will be tried over remote targets)
		ver, _ = semver.ParseVersion(selector)
	}
	// check whether requested version is already installed
	if ver != nil {
		local, err := Ls()
		if err != nil {
			return "", err
		}
		for _, v := range local {
			if ver.Equals(v) {
				return ver.String(), nil
			}
		}
	}
	// ... apparently it's not
	if releaseMap == nil {
		ver = nil
		rng, err := semver.ParseRange(selector)
		if err != nil {
			return "", err
		}
		releaseMap, err = LsRemote()
		if err != nil {
			return "", err
		}
		var vs = make([]*semver.Version, len(releaseMap))
		var i = 0
		for k := range releaseMap {
			vs[i] = k
			i++
		}
		sort.Sort(sort.Reverse(semver.VersionSlice(vs)))
		for _, v := range vs {
			if rng.Contains(v) {
				ver = v
				break
			}
		}
		if ver == nil {
			tt := make([]string, len(vs))
			for i, v := range vs {
				tt[i] = v.String()
			}
			return "", errors.New("No compatible version found for " + selector +
			"\nValid install targets: " + strings.Join(tt, ", "))
		}
	}
	url := releaseMap[ver]
	if matched, _ := regexp.MatchString("^\\w+[+]\\w+://", url); !matched {
		return "", errors.New("URL must contain qualifier, e.g. tgz+http://...")
	}
	var fileType string = url[0:strings.Index(url, "+")]
	url = url[strings.Index(url, "+") + 1:]
	var file string
	var deleteFileWhenFinnished bool
	if strings.HasPrefix(url, "file://") {
		file = strings.TrimPrefix(url, "file://")
		if runtime.GOOS == "windows" {
			// file:///C:/path/...
			file = strings.Replace(strings.TrimPrefix(file, "/"), "/", "\\", -1)
		}
	} else {
		log.Info("Downloading ", ver, " (", url, ")")
		file, err = download(url, fileType)
		if err != nil {
			return "", err
		}
		deleteFileWhenFinnished = true
	}
	switch runtime.GOOS {
	case "darwin":
		err = installOnDarwin(ver.String(), file, fileType)
	case "linux":
		err = installOnLinux(ver.String(), file, fileType)
	case "windows":
		err = installOnWindows(ver.String(), file, fileType)
	default:
		err = errors.New(runtime.GOOS + " OS is not supported")
	}
	if err == nil && deleteFileWhenFinnished {
		os.Remove(file)
	}
	return ver.String(), err
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

func download(url string, fileType string) (file string, err error) {
	tmp, err := ioutil.TempFile("", "jabba-d-")
	if err != nil {
		return
	}
	if fileType == "exe" {
		err = tmp.Close()
		if err != nil {
			return
		}
		err = os.Rename(tmp.Name(), tmp.Name() + ".exe")
		if err != nil {
			return
		}
		tmp, err = os.OpenFile(tmp.Name() + ".exe", os.O_RDWR, 0600)
		if err != nil {
			return
		}
	}
	defer tmp.Close()
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
	if strings.Contains(url, "zulu") {
		req.Header.Set("Referer", "http://www.azul.com/downloads/zulu/")
	}
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

func installOnDarwin(ver string, file string, fileType string) (err error) {
	target := cfg.Dir() + "/jdk/" + ver
	switch fileType {
	case "dmg":
		err = installFromDmg(file, target)
	case "zip":
		err = installFromZip(file, target + "/Contents/Home")
	default:
		return errors.New(fileType + " is not supported")
	}
	if err == nil {
		err = assertJavaDistribution(target)
	}
	if err != nil {
		os.RemoveAll(target)
	}
	return
}

func installFromDmg(source string, target string) error {
	tmp, err := ioutil.TempDir("", "jabba-i-")
	if err != nil {
		return err
	}
	basename := filepath.Base(source)
	mountpoint := tmp + "/" + basename
	pkgdir := tmp + "/" + basename + "-pkg"
	err = executeInShell([][]string{
		[]string{"Mounting " + source, "hdiutil mount -mountpoint " + mountpoint + " " + source},
		[]string{"Extracting " + source + " to " + target,
			"pkgutil --expand " + mountpoint + "/*.pkg " + pkgdir},
		[]string{"", "mkdir -p " + target},

		// todo: instead of relying on a certain pkg structure - find'n'extract all **/*/Payload

		// oracle
		[]string{"",
			"if [ -f " + pkgdir + "/jdk*.pkg/Payload" + " ]; then " +
			"cd " + pkgdir + "/jdk*.pkg && " +
			"cat Payload | gzip -d | cpio -i && " +
			"mv Contents " + target + "/" +
			"; fi"},

		// apple
		[]string{"",
			"if [ -f " + pkgdir + "/JavaForOSX.pkg/Payload" + " ]; then " +
			"cd " + pkgdir + "/JavaForOSX.pkg && " +
			"cat Payload | gzip -d | cpio -i && " +
			"mv Library/Java/JavaVirtualMachines/*/Contents " + target + "/" +
			"; fi"},

		[]string{"Unmounting " + source, "hdiutil unmount " + mountpoint},
	})
	if err == nil {
		os.RemoveAll(tmp)
	}
	return err
}

func installOnLinux(ver string, file string, fileType string) (err error) {
	target := cfg.Dir() + "/jdk/" + ver
	switch fileType {
	case "bin":
		err = installFromBin(file, target)
	case "ia":
		err = installFromIa(file, target)
	case "tgz":
		err = installFromTgz(file, target)
	case "zip":
		err = installFromZip(file, target)
	default:
		return errors.New(fileType + " is not supported")
	}
	if err == nil {
		err = assertJavaDistribution(target)
	}
	if err != nil {
		os.RemoveAll(target)
	}
	return
}

func installOnWindows(ver string, file string, fileType string) (err error) {
	target := filepath.Join(cfg.Dir(), "jdk", ver)
	switch fileType {
	case "exe":
		err = installFromExe(file, target)
	case "zip":
		err = installFromZip(file, target)
	default:
		return errors.New(fileType + " is not supported")
	}
	if err == nil {
		err = assertJavaDistribution(target)
	}
	if err != nil {
		os.RemoveAll(target)
	}
	return
}

func installFromBin(source string, target string) (err error) {
	tmp, err := ioutil.TempDir("", "jabba-i-")
	if err != nil {
		return
	}
	err = executeInShell([][]string{
		[]string{"", "cp " + source + " " + tmp},
		[]string{"Extracting " + filepath.Join(tmp, filepath.Base(source)) + " to " + target,
			"cd " + tmp + " && echo | sh " + filepath.Base(source) + " && mv jdk*/ " + target},
	})
	if err == nil {
		os.RemoveAll(tmp)
	}
	return
}

func installFromIa(source string, target string) (err error) {
	tmp, err := ioutil.TempDir("", "jabba-i-")
	if err != nil {
		return
	}
	err = executeInShell([][]string{
		[]string{"", "printf 'LICENSE_ACCEPTED=TRUE\\nUSER_INSTALL_DIR=" + target + "' > " +
			filepath.Join(tmp, "installer.properties")},
		[]string{"Extracting " + source + " to " + target,
			"echo | sh " + source + " -i silent -f " + filepath.Join(tmp, "installer.properties")},
	})
	if err == nil {
		os.RemoveAll(tmp)
	}
	return
}

func installFromExe(source string, target string) error {
	log.Info("Unpacking " + source + " to " + target)
	// using ShellExecute instead of exec.Command so user could decide whether to trust the installer when UAC is active
	return w32.ShellExecuteAndWait(w32.HWND(0), "open", source, "/s INSTALLDIR=\"" + target +
		"\" STATIC=1 AUTO_UPDATE=0 WEB_JAVA=0 WEB_ANALYTICS=0 REBOOT=0", "", 3)
}

func installFromTgz(source string, target string) error {
	return executeInShell([][]string{
		[]string{"", "mkdir -p " + target},
		[]string{"Extracting " + source + " to " + target,
			"tar xzf " + source + " --strip-components=1 -C " + target},
	})
}

func installFromZip(source string, target string) error {
	log.Info("Extracting " + source + " to " + target)
	return unzip(source, target, true)
}

func unzip(source string, target string, strip bool) error {
	r, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer r.Close()
	var prefixToStrip = ""
	if strip {
		entriesPerLevel := make(map[int]int)
		prefixMap := make(map[int]string)
		for _, f := range r.File {
			level := 0
			for _, c := range f.Name {
				if c == '/' {
					level++
				}
			}
			if !f.Mode().IsDir() {
				level++
			} else {
				prefixMap[level] = f.Name
			}
			entriesPerLevel[level]++
		}
		for i := 0; i < len(entriesPerLevel); i++ {
			if entriesPerLevel[i] > 1 && i > 0 {
				prefixToStrip = prefixMap[i - 1]
				break
			}
		}
	}
	for _, f := range r.File {
		name := strings.TrimPrefix(f.Name, prefixToStrip)
		if f.Mode().IsDir() {
			os.MkdirAll(filepath.Join(target, name), 0755)
		} else {
			fr, err := f.Open()
			if err != nil {
				return err
			}
			f, err := os.OpenFile(filepath.Join(target, name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(f, fr)
			if err != nil {
				return err
			}
			f.Close()
		}
	}
	return nil
}

func executeInShell(cmd [][]string) error {
	for _, command := range cmd {
		if command[0] != "" {
			log.Info(command[0])
		}
		var execArg []string
		if runtime.GOOS == "windows" {
			execArg = []string{"cmd", "/C"}
		} else {
			execArg = []string{"sh", "-c"}
		}
		out, err := exec.Command(execArg[0], execArg[1], command[1]).CombinedOutput()
		if err != nil {
			log.Error(string(out))
			return errors.New("'" + command[1] + "' failed: " + err.Error())
		}
	}
	return nil
}

func assertJavaDistribution(target string) error {
	if runtime.GOOS == "darwin" {
		target += "/Contents/Home"
	}
	var javaBin = "java"
	if runtime.GOOS == "windows" {
		javaBin += ".exe"
	}
	var path = filepath.Join(target, "bin", javaBin)
	var err error
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = errors.New(path + " wasn't found. " +
		"If you believe this is an error - please create a ticket at https://github.com/shyiko/jabba/issue " +
		"(specify OS and command that was used)")
	}
	return err
}
