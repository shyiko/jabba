package semver

import (
	"github.com/wmark/semver"
	"strings"
	"fmt"
)

type Range struct {
	qualifier string
	raw       string
	rng       *semver.Range
}

func (l *Range) Contains(r *Version) bool {
	return l.qualifier == r.qualifier && l.rng.Contains(r.ver)
}

func (t *Range) String() string {
	return t.raw
}

func ParseRange(raw string) (*Range, error) {
	p := new(Range)
	p.raw = raw
	// selector can be either <version> or <qualifier>@<version>
	if strings.Contains(raw, "@") {
		p.qualifier = raw[0:strings.Index(raw, "@")]
		raw = raw[strings.Index(raw, "@") + 1:len(raw)]
	}
	parsed, err := semver.NewRange(raw)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid version", raw)
	}
	p.rng = parsed
	return p, nil
}
