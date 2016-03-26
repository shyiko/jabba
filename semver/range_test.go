package semver

import (
	"testing"
)

func TestContains(t *testing.T) {
	assertWithinRange(t, "1.8", "1.8.0", true)
	assertWithinRange(t, "1.7", "1.8.0", false)
	assertWithinRange(t, "1.8.0-0", "1.8.0-0", true)
	assertWithinRange(t, "1.8.0-0", "1.8.0-1", false)
	assertWithinRange(t, "~1.8", "1.8.99", true)
	assertWithinRange(t, "~1.8", "1.9.0", false)
	assertWithinRange(t, "a@1.8", "a@1.8.72", true)
	assertWithinRange(t, "1.8", "a@1.8.72", false)
	assertWithinRange(t, "a@1.8", "b@1.8.72", false)
	assertWithinRange(t, "a@1.8", "1.8.72", false)
}

func assertWithinRange(t *testing.T, rng string, ver string, value bool) {
	r, err := ParseRange(rng)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	v, err := ParseVersion(ver)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if r.Contains(v) != value {
		t.Fatalf("expected range %v to contain %v (%v)", rng, ver, value)
	}
}
