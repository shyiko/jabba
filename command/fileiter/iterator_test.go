package fileiter

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestDFS(t *testing.T) {
	ok := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	dir, err := ioutil.TempDir("", "fileiter_test")
	ok(err)
	ok(touch(dir, "a", "b", "c"))
	ok(touch(dir, "b", "c"))
	ok(touch(dir, "c"))
	ok(mkdir(dir, "d"))
	expectedSeq := []string{
		"a",
		filepath.Join("a", "b"),
		filepath.Join("a", "b", "c"),
		"b",
		filepath.Join("b", "c"),
		"c",
		"d",
	}
	test(t, dir, expectedSeq)
	test(t, filepath.Join(dir, "d"), nil)
	test(t, filepath.Join(dir, "non-existent"), nil)
}

func TestBFS(t *testing.T) {
	ok := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	dir, err := ioutil.TempDir("", "fileiter_test")
	ok(err)
	ok(touch(dir, "a", "b", "c"))
	ok(touch(dir, "b", "c"))
	ok(touch(dir, "c"))
	ok(mkdir(dir, "d"))
	expectedSeq := []string{
		"a",
		"b",
		"c",
		"d",
		filepath.Join("a", "b"),
		filepath.Join("b", "c"),
		filepath.Join("a", "b", "c"),
	}
	test(t, dir, expectedSeq, BreadthFirst())
	test(t, filepath.Join(dir, "d"), nil, BreadthFirst())
	test(t, filepath.Join(dir, "non-existent"), nil, BreadthFirst())
}

func TestNew(t *testing.T) {
	ok := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	dir, err := ioutil.TempDir("", "fileiter_test")
	ok(err)
	ok(touch(dir, "a", "b"))
	if it := New(filepath.Join(dir, "a", "b")); it.Err() != nil || it.IsDir() {
		t.Fatal(it.Err())
	}
	if it := New(filepath.Join(dir, "a")); it.Err() != nil || !it.IsDir() {
		t.Fatal(it.Err())
	}
	if it := New(filepath.Join(dir, "b")); it.Err() == nil {
		t.Fatal()
	}
}

func test(t *testing.T, dir string, expectedSeq []string, opts ...IterationOption) {
	i := 0
	for w := New(dir, opts...); w.Next(); {
		if w.Err() != nil {
			t.Fatal(w.Err())
		}
		expected := filepath.Join(dir, expectedSeq[i])
		actual := filepath.Join(w.Dir(), w.Name())
		if actual != expected {
			t.Fatalf("actual: %v != expected: %v", actual, expected)
		}
		i++
	}
}

func touch(path ...string) error {
	filename := filepath.Join(path...)
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, nil, 0755); err != nil {
		return err
	}
	return nil
}

func mkdir(path ...string) error {
	if err := os.MkdirAll(filepath.Join(path...), 0755); err != nil {
		return err
	}
	return nil
}
