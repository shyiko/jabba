package semver

import (
	"fmt"
	"github.com/Masterminds/semver"
	"strings"
)

type Version struct {
	qualifier string
	raw       string
	ver       *semver.Version
}

func (l *Version) LessThan(r *Version) bool {
	if l.qualifier == r.qualifier {
		return l.ver.LessThan(r.ver)
	}
	return l.qualifier > r.qualifier
}

func (l *Version) Equals(r *Version) bool {
	return l.raw == r.raw
}

func (t *Version) String() string {
	return t.raw
}

func ParseVersion(raw string) (*Version, error) {
	p := new(Version)
	p.raw = raw
	// selector can be either <version> or <qualifier>@<version>
	if strings.Contains(raw, "@") {
		p.qualifier = raw[0:strings.Index(raw, "@")]
		raw = raw[strings.Index(raw, "@")+1:]
	}
	parsed, err := semver.NewVersion(raw)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid version", raw)
	}
	p.ver = parsed
	return p, nil
}

type VersionSlice []*Version

// impl sort.Interface:

func (c VersionSlice) Len() int {
	return len(c)
}
func (c VersionSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c VersionSlice) Less(i, j int) bool {
	return c[i].LessThan(c[j])
}
