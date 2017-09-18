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
	assertWithinRange(t, "a@1.7.x", "a@1.7.72", true)
	assertWithinRange(t, "a@1.7.x", "a@1.8.72", false)
	assertWithinRange(t, "a@>=1.7 <=1.8.75", "a@1.8.72", true)
	assertWithinRange(t, "a@>=1.7 <=1.8.75", "a@1.8.80", false)
}

func TestPre070Compat(t *testing.T) {
	actual := pre070Compat("1.1, >= 1.2 <= 1.3.0, 4, 1.5,1.6.0, ~1.2")
	expected := "1.1.x, >= 1.2, <= 1.3.0, 4.x, 1.5.x,1.6.0, ~1.2"
	if actual != expected {
		t.Fatalf("actual: %#v != expected: %#v", actual, expected)
	}
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
