package command

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestBinJavaRelocation(t *testing.T) {
	ok := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	nok := func(err error) {
		if err == nil {
			t.Fatal(err)
		}
	}
	dir, err := ioutil.TempDir("", "install_test")
	ok(err)
	for _, scenario := range []struct {
		os     string
		bin    string
		prefix string
		paths  []string
	}{
		{
			os:     "linux",
			bin:    "java",
			prefix: "",
			paths:  []string{""},
		},
		{
			os:     "darwin",
			bin:    "java",
			prefix: filepath.Join("Contents", "Home"),
			paths: []string{
				"",
				filepath.Join("Home"),
				filepath.Join("Contents", "Home"),
			},
		},
		{
			os:     "windows",
			bin:    "java.exe",
			prefix: "",
			paths:  []string{""},
		},
	} {
		for _, p := range scenario.paths {
			test1 := filepath.Join(dir, "test1")
			ok(touch(test1, p, "bin", scenario.bin))
			ok(normalizePathToBinJava(test1, scenario.os))
			ok(file(test1, scenario.prefix, "bin", scenario.bin))

			test2 := filepath.Join(dir, "test2")
			ok(touch(test2, "subdir", p, "bin", scenario.bin))
			ok(normalizePathToBinJava(test2, scenario.os))
			ok(file(test2, scenario.prefix, "bin", scenario.bin))

			test3 := filepath.Join(dir, "test3")
			ok(touch(test3, "subdir", "subdir", p, "bin", scenario.bin))
			ok(normalizePathToBinJava(test3, scenario.os))
			ok(file(test3, scenario.prefix, "bin", scenario.bin))

			test4 := filepath.Join(dir, "test4")
			ok(touch(test4, "file"))
			ok(touch(test4, "subdir", "subdir", p, "bin", scenario.bin))
			ok(normalizePathToBinJava(test4, scenario.os))
			ok(file(test4, scenario.prefix, "bin", scenario.bin))

			test5 := filepath.Join(dir, "test5")
			ok(touch(test5, "bin", "file"))
			nok(normalizePathToBinJava(test5, scenario.os))
			ok(file(test5, "bin", "file"))
		}
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

func file(path ...string) error {
	if _, err := os.Stat(filepath.Join(path...)); os.IsNotExist(err) {
		return err
	}
	return nil
}
