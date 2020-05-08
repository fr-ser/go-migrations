package utils

import (
	"io/ioutil"
	"path/filepath"
	"syscall"
	"testing"
)

func TestFileExistsFalse(t *testing.T) {

	if result, _ := FileExists(filepath.FromSlash("/hell/if/this/ever/exists")); result {
		t.Errorf("Returned true for non existing file ")
	}
}

func TestFileExistsTrue(t *testing.T) {
	f, _ := ioutil.TempFile("", "tmp_file")
	defer syscall.Unlink(f.Name())

	if result, _ := FileExists(filepath.FromSlash(f.Name())); !result {
		t.Errorf("Returned false for an existing file ")
	}
	if result, _ := FileExists(filepath.Dir(filepath.FromSlash(f.Name()))); !result {
		t.Errorf("Returned false for an existing file ")
	}
}
