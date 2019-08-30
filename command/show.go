package command

import (
	"fmt"
	"net/url"
	"runtime"
	"strings"

	"github.com/shyiko/jabba/semver"
	"github.com/spf13/pflag"
)

const notAvailable = "Not Available"

type versionInfo struct {
	name             string
	isInstalled      bool
	downloadURL      string
	javaMajorVersion string
	vendor           string
	err              error
}

func newVeriosnInfo(selector string) *versionInfo {
	return &versionInfo{name: selector,
		isInstalled:      false,
		downloadURL:      notAvailable,
		vendor:           notAvailable,
		javaMajorVersion: notAvailable,
	}
}

func (v *versionInfo) formatString() string {
	templ := "%s\n" +
		"  Installed:\t\t%s\n" +
		"  DownloadUrl:\t\t%s\n" +
		"  JavaMajorVersion:\t%s\n" +
		"  Vendor:\t\t%s\n"
	isInstalled := "No"
	if v.isInstalled {
		isInstalled = "Yes"
	}
	return fmt.Sprintf(templ, v.name, isInstalled, v.downloadURL, v.javaMajorVersion, v.vendor)
}

// Show prints the information about given slelectors
func Show(args []string) error {
	if len(args) == 0 {
		return pflag.ErrHelp
	}
	for _, s := range args {
		info := getVersionInfoFor(s)
		if info.err != nil {
			fmt.Printf("%s: %s\n", s, info.err)
		} else {
			fmt.Println(info.formatString())
		}
	}

	return nil
}

// This is a wrapper for Ls which allow to mock Ls functionality for tests
var ls = func() ([]*semver.Version, error) { return Ls() }

// This is a wrapper for Ls which allow to mock LsRemote functionality for tests
var lsRemote = func() (map[*semver.Version]string, error) {
	return LsRemote(runtime.GOOS, runtime.GOARCH)
}

func getVersionInfoFor(selector string) *versionInfo {
	ver, err := semver.ParseVersion(selector)
	if err != nil {
		return &versionInfo{err: fmt.Errorf("%s is not a valid version", selector)}
	}
	info := newVeriosnInfo(selector)
	installed, err := ls()
	for _, v := range installed {
		if v.Equals(ver) {
			info.isInstalled = true
			info.javaMajorVersion = fmt.Sprintf("%d.%d", v.Major(), v.Minor())
		}
	}
	releaseMap, _ := lsRemote()
	for v, url := range releaseMap {
		if v.Equals(ver) {
			info.downloadURL = strings.Join(strings.Split(url, "+")[1:], "")
			info.vendor = getVendor(info.downloadURL)
			info.javaMajorVersion = fmt.Sprintf("%d.%d", v.Major(), v.Minor())
		}
	}
	return info
}

func getVendor(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	parts := strings.Split(u.Hostname(), ".")
	fmt.Println(parts)
	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
	return domain
}
