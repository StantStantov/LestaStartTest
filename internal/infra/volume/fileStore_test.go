//go:build integration || !unit

package volume_test

import (
	"Stant/LestaGamesInternship/internal/app/config"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"Stant/LestaGamesInternship/internal/infra/volume"
	"crypto/rand"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileStore(t *testing.T) {
	dirpath := filepath.Join(os.Getenv(config.PathToDocsEnv), "FileStore")
	if err := os.RemoveAll(dirpath); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dirpath, 0770); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dirpath); err != nil {
		t.Fatal(err)
	}
	fileStore := volume.NewFileStore(dirpath)

	t.Run("Test Save File", func(t *testing.T) {
		t.Parallel()

		testSaveFile(t, fileStore)
	})
	t.Run("Test Open File", func(t *testing.T) {
		t.Parallel()

		testOpenFile(t, fileStore)
	})
	t.Run("Test Rename File", func(t *testing.T) {
		t.Parallel()

		testRenameFile(t, fileStore)
	})
	t.Run("Test Delete File", func(t *testing.T) {
		t.Parallel()

		testDeleteFile(t, fileStore)
	})
}

func testSaveFile(t *testing.T, fileStore stores.FileStore) {
	t.Helper()

	t.Run("PASS if saved", func(t *testing.T) {
		t.Parallel()

		want := true
		filename := rand.Text()
		data := strings.NewReader("Hello, world!")

		if err := fileStore.Save(filename, data); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		got, err := fileStore.IsExist(filename)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if want != got {
			t.Fatalf("Wanted %v, got %v", want, got)
		}
	})
}

func testOpenFile(t *testing.T, fileStore stores.FileStore) {
	t.Helper()

	t.Run("PASS if opened", func(t *testing.T) {
		t.Parallel()

		stringBuilder := new(strings.Builder)
		wantName := rand.Text()
		wantData := "Hello, world!"

		if err := fileStore.Save(wantName, strings.NewReader(wantData)); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		file, err := fileStore.Open(wantName)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		fi, err := file.Stat()
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		stringBuilder.Grow(int(fi.Size()))
		io.Copy(stringBuilder, file)

		gotName := file.Name()
		gotData := stringBuilder.String()
		if wantName != gotName {
			t.Fatalf("Wanted %v, got %v", wantName, gotName)
		}
		if wantData != gotData {
			t.Fatalf("Wanted %v, got %v", wantData, gotData)
		}
	})
	t.Run("FAIL if doesn't exists", func(t *testing.T) {
		t.Parallel()

		filename := rand.Text()
		if _, err := fileStore.Open(filename); err == nil {
			t.Fatalf("Wanted err, got %v", err)
		}
	})
}

func testRenameFile(t *testing.T, fileStore stores.FileStore) {
	t.Helper()

	t.Run("PASS if renamed", func(t *testing.T) {
		t.Parallel()

		wantName := rand.Text()
		filename := rand.Text()
		data := "Hello, world!"

		if err := fileStore.Save(filename, strings.NewReader(data)); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := fileStore.Rename(filename, wantName); err != nil {
			t.Errorf("Wanted %v, got %v", nil, err)
		}
	})
	t.Run("FAIL if doesn't exists", func(t *testing.T) {
		t.Parallel()

		oldName := rand.Text()
		newName := rand.Text()
		if err := fileStore.Rename(oldName, newName); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}

func testDeleteFile(t *testing.T, fileStore stores.FileStore) {
	t.Helper()

	t.Run("PASS if deleted", func(t *testing.T) {
		t.Parallel()

		wantName := rand.Text()
		data := "Hello, world!"

		if err := fileStore.Save(wantName, strings.NewReader(data)); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := fileStore.Delete(wantName); err != nil {
			t.Errorf("Wanted %v, got %v", nil, err)
		}
	})
	t.Run("FAIL if doesn't exists", func(t *testing.T) {
		t.Parallel()

		filename := rand.Text()
		if err := fileStore.Delete(filename); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}
