package command

import (
	"github.com/shyiko/jabba/cfg"
	"os"
	"reflect"
	"testing"
)

func TestDeactivate(t *testing.T) {
	prevPath := os.Getenv("PATH")
	defer func() { os.Setenv("PATH", prevPath) }()
	os.Setenv("PATH", "/usr/local/bin:"+cfg.Dir()+"/jdk/zulu@1.8.72/bin:/system-jdk/bin:/usr/bin")
	os.Setenv("JAVA_HOME", cfg.Dir()+"/jdk/zulu@1.8.72")
	os.Setenv("JAVA_HOME_BEFORE_JABBA", "/system-jdk")
	actual, err := Deactivate()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected := []string{
		"export PATH=\"/usr/local/bin:/system-jdk/bin:/usr/bin\"",
		"export JAVA_HOME=\"/system-jdk\"",
		"unset JAVA_HOME_BEFORE_JABBA",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual: %v != expected: %v", actual, expected)
	}
}

func TestDeactivateInUnusedEnv(t *testing.T) {
	prevPath := os.Getenv("PATH")
	defer func() { os.Setenv("PATH", prevPath) }()
	os.Setenv("PATH", "/usr/local/bin:/system-jdk/bin:/usr/bin")
	os.Setenv("JAVA_HOME", "/system-jdk")
	os.Unsetenv("JAVA_HOME_BEFORE_JABBA")
	actual, err := Deactivate()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected := []string{
		"export PATH=\"/usr/local/bin:/system-jdk/bin:/usr/bin\"",
		"export JAVA_HOME=\"/system-jdk\"",
		"unset JAVA_HOME_BEFORE_JABBA",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual: %v != expected: %v", actual, expected)
	}
}
