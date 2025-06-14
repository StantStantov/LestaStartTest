package apptest

import (
	"os"
	"path/filepath"
	"testing"
)

func CreateTestDir(t *testing.T, dirpath, dirname string) string {
	t.Helper()

	dirpath = filepath.Join(dirpath, dirname)
	if err := os.RemoveAll(dirpath); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dirpath, 0770); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dirpath); err != nil {
		t.Fatal(err)
	}

	return dirpath
}

func CreateTestFile(t *testing.T, prefix string) *os.File {
	t.Helper()

	file, err := os.CreateTemp("", prefix+"_*")
	if err != nil {
		t.Fatalf("Wanted %v, got %v", nil, err)
	}

	t.Cleanup(func() {
		file.Close()
		os.Remove(file.Name())
	})

	return file
}
