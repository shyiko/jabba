package semver

import (
	"testing"
	"reflect"
)

func TestSort(t *testing.T) {
	actual := Sort([]string{"0.2.0", "0.1.20", "0.1.10", "0.1.2"})
	expected := []string{"0.2.0", "0.1.20", "0.1.10", "0.1.2"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual: %v != expected: %v", actual, expected)
	}
}
