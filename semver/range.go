package semver

import (
	"fmt"
	"github.com/Masterminds/semver"
	"regexp"
	"strings"
)

var pre070CompatRegexp = regexp.MustCompile("(^|,\\s*)\\d+([.]\\d+)?[.]?")
var pre070CompatRepl = func(input string) string {
	if strings.HasSuffix(input, ".") {
		return input
	}
	return input + ".x"
}

func pre070Compat(version string) string {
	// 1.2 -> 1.2.x
	s := pre070CompatRegexp.ReplaceAllStringFunc(version, pre070CompatRepl)
	// >= 1.2 <= 2.4 -> >= 1.2, <= 2.4
	split := regexp.MustCompile("\\s+").Split(s, -1)
	for i, v := range split {
		if i > 0 {
			switch v[0] {
			case '!', '=', '<', '>', '~', '^':
				pv := split[i-1]
				pvc := pv[len(pv)-1]
				if pvc != '|' && pvc != ',' {
					split[i-1] = pv + ","
				}
			}
		}
	}
	return strings.Join(split, " ")
}

type Range struct {
	qualifier string
	raw       string
	rng       *semver.Constraints
}

func (l *Range) Contains(r *Version) bool {
	return l.qualifier == r.qualifier && l.rng.Check(r.ver)
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
		raw = raw[strings.Index(raw, "@")+1:]
		if raw == "" {
			// `jabba ls-remote zulu@` convenience
			raw = ">=0.0.0-0"
		}
	}
	constraint := pre070Compat(raw)
	parsed, err := semver.NewConstraint(constraint)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid version", p.raw)
	}
	p.rng = parsed
	return p, nil
}
