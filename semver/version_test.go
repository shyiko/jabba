package semver

import (
	"testing"
	"reflect"
	"sort"
)

func TestSort(t *testing.T) {
	actual := asVersionSlice(t,
		"0.2.0", "a@1.8.10", "b@1.8.2", "0.1.20", "a@1.8.2", "0.1.10", "0.1.2")
	sort.Sort(sort.Reverse(VersionSlice(actual)))
	expected := asVersionSlice(t,
		"0.2.0", "0.1.20", "0.1.10", "0.1.2", "a@1.8.10", "a@1.8.2", "b@1.8.2")
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual: %v != expected: %v", actual, expected)
	}
}

func asVersionSlice(t *testing.T, slice ...string) (r []*Version) {
	for _, value := range slice {
		ver, err := ParseVersion(value)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		r = append(r, ver)
	}
	return
}
