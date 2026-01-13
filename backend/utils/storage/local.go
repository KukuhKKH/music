package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	Path string
}

func NewLocalStorage(path string) (*LocalStorage, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return nil, err
		}
	}
	return &LocalStorage{Path: path}, nil
}

func (s *LocalStorage) Upload(ctx context.Context, filename string, file io.Reader) (string, error) {
	dstPath := filepath.Join(s.Path, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (s *LocalStorage) Delete(filename string) error {
	return os.Remove(filepath.Join(s.Path, filename))
}

func (s *LocalStorage) GetURL(filename string) string {
	return fmt.Sprintf("http://localhost:8000/storage/%s", filename)
}
