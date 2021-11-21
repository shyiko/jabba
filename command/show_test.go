package command

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/shyiko/jabba/semver"
)

func TestGetVersionInfoForGivenSelectorForInvalidVersion(t *testing.T) {
	selector := "sdfsd"
	info := getVersionInfoFor(selector)

	if info.err == nil {
		t.Errorf("Expected error. Got %s", info.err.Error())
	}
	expectedError := fmt.Sprintf("%s is not a valid version", selector)
	if info.err.Error() != expectedError {
		t.Errorf("Expected: %s\nGot: %s", expectedError, info.err.Error())
	}
}

func TestGetVersionInfoForGivenSelectorForValidVersionButNotAvailableAndNotInstalled(t *testing.T) {
	selector := "openjdk@1.1.1"
	ls = func() ([]*semver.Version, error) { return []*semver.Version{}, nil }
	lsRemote = func() (map[*semver.Version]string, error) { return make(map[*semver.Version]string, 0), nil }
	info := getVersionInfoFor(selector)

	if info.err != nil {
		t.Errorf("Expected error to be nil. Got %s", info.err.Error())
	}
	expected := &versionInfo{
		name:             selector,
		downloadURL:      notAvailable,
		isInstalled:      false,
		javaMajorVersion: notAvailable,
		vendor:           notAvailable,
		err:              nil,
	}
	if !reflect.DeepEqual(info, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, info)
	}
}

func TestGetVersionInfoForGivenSelectorForValidVersionInstalled(t *testing.T) {
	selector := "openjdk@1.1.1"
	v, _ := semver.ParseVersion(selector)
	ls = func() ([]*semver.Version, error) {
		return []*semver.Version{v}, nil
	}
	lsRemote = func() (map[*semver.Version]string, error) { return make(map[*semver.Version]string, 0), nil }
	info := getVersionInfoFor(selector)

	if info.err != nil {
		t.Errorf("Expected error to be nil. Got %s", info.err.Error())
	}
	expected := &versionInfo{
		name:             selector,
		downloadURL:      notAvailable,
		isInstalled:      true,
		javaMajorVersion: "1.1",
		vendor:           notAvailable,
		err:              nil,
	}
	if !reflect.DeepEqual(info, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, info)
	}
}

func TestGetVersionInfoForGivenSelectorForValidVersionAvailableButNotIsnstall(t *testing.T) {
	selector := "openjdk@1.1.1"
	v, _ := semver.ParseVersion(selector)
	ls = func() ([]*semver.Version, error) { return []*semver.Version{}, nil }
	releaseMap := map[*semver.Version]string{}
	releaseMap[v] = "tgz+https://dl.foo.com/some/some-java.tar.gz"
	lsRemote = func() (map[*semver.Version]string, error) { return releaseMap, nil }
	info := getVersionInfoFor(selector)

	if info.err != nil {
		t.Errorf("Expected error to be nil. Got %s", info.err.Error())
	}
	expected := &versionInfo{
		name:             selector,
		downloadURL:      "https://dl.foo.com/some/some-java.tar.gz",
		isInstalled:      false,
		javaMajorVersion: "1.1",
		vendor:           "foo.com",
		err:              nil,
	}
	if !reflect.DeepEqual(info, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, info)
	}
}

func TestGetVersionInfoForGivenSelectorForValidVersionAvailableAndIsnstalled(t *testing.T) {
	selector := "openjdk@1.1.1"
	v, _ := semver.ParseVersion(selector)
	ls = func() ([]*semver.Version, error) { return []*semver.Version{v}, nil }
	releaseMap := map[*semver.Version]string{}
	releaseMap[v] = "tgz+https://dl.foo.com/some/some-java.tar.gz"
	lsRemote = func() (map[*semver.Version]string, error) { return releaseMap, nil }
	info := getVersionInfoFor(selector)

	if info.err != nil {
		t.Errorf("Expected error to be nil. Got %s", info.err.Error())
	}
	expected := &versionInfo{
		name:             selector,
		downloadURL:      "https://dl.foo.com/some/some-java.tar.gz",
		isInstalled:      true,
		javaMajorVersion: "1.1",
		vendor:           "foo.com",
		err:              nil,
	}
	if !reflect.DeepEqual(info, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, info)
	}
}

func TestGetVersionInfoForGivenSelectorFormatString(t *testing.T) {
	selector := "openjdk@1.1.1"
	v, _ := semver.ParseVersion(selector)
	ls = func() ([]*semver.Version, error) { return []*semver.Version{v}, nil }
	releaseMap := map[*semver.Version]string{}
	releaseMap[v] = "tgz+https://dl.foo.com/some/some-java.tar.gz"
	lsRemote = func() (map[*semver.Version]string, error) { return releaseMap, nil }
	info := getVersionInfoFor(selector)

	if info.err != nil {
		t.Errorf("Expected error to be nil. Got %s", info.err.Error())
	}
	expected := `openjdk@1.1.1
  Installed:		Yes
  DownloadUrl:		https://dl.foo.com/some/some-java.tar.gz
  JavaMajorVersion:	1.1
  Vendor:		foo.com
`
	if expected != info.formatString() {
		t.Errorf("\nExpected: %v\nGot: %v", expected, info.formatString())
	}
}
