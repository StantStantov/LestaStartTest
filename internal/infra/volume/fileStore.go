package volume

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type FileStore struct {
	dirpath string
}

func NewFileStore(dirpath string) *FileStore {
	return &FileStore{}
}

func (s *FileStore) Save(filename string, data io.Reader) error {
	file, err := os.Create(filepath.Join(s.dirpath, filename))
	if err != nil {
		return fmt.Errorf("volume/fileStore.Save: [%w]", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, data); err != nil {
		return fmt.Errorf("volume/fileStore.Save: [%w]", err)
	}

	return nil
}

func (s *FileStore) IsExist(filename string) (bool, error) {
	info, err := os.Stat(filepath.Join(s.dirpath, filename))
	if errors.Is(err, fs.ErrNotExist) {
		return false, fmt.Errorf("volume/fileStore.IsExist: [%w]", err)
	}
	if info == nil  {
		return false, fmt.Errorf("volume/fileStore.IsExist: [Couldn't Stat info about file %q]", filename)
	}
	return !info.IsDir(), nil
}

func (s *FileStore) Open(filename string) (*os.File, error) {
	file, err := os.Open(filepath.Join(s.dirpath, filename))
	if err != nil {
		return nil, fmt.Errorf("volume/fileStore.Open: [%w]", err)
	}

	return file, nil
}

func (s *FileStore) Rename(oldName, newName string) error {
	oldPath := filepath.Join(s.dirpath, oldName)
	newPath := filepath.Join(s.dirpath, newName)
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("volume/fileStore.Rename: [%w]", err)
	}

	return nil
}

func (s *FileStore) Delete(filename string) error {
	if err := os.Remove(filepath.Join(s.dirpath, filename)); err != nil {
		return fmt.Errorf("volume/fileStore.Delete: [%w]", err)
	}

	return nil
}
