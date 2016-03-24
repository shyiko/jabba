package semver

import (
	"sort"
	"github.com/wmark/semver"
)

type Version struct {
	Raw    string
	Parsed *semver.Version
}

func NewVersion(raw string) *Version {
	p := new(Version)
	p.Raw = raw
	parsed, err := semver.NewVersion(raw)
	if err != nil {
		panic(raw + " is not a valid version")
	}
	p.Parsed = parsed
	return p
}

type VersionSlice []*Version

func (c VersionSlice) Len() int {
	return len(c)
}
func (c VersionSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c VersionSlice) Less(i, j int) bool {
	return c[i].Parsed.Less(c[j].Parsed)
}

func Sort(vs []string) []string {
	var svs = make([]*Version, len(vs))
	for i, v := range vs {
		svs[i] = NewVersion(v)
	}
	sort.Sort(sort.Reverse(VersionSlice(svs)))
	for i, v := range svs {
		vs[i] = v.Raw
	}
	return vs
}

