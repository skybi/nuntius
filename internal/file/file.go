package file

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Write writes bytes to a file, creating the file if it does not exist and overriding it if it does
func Write(path string, contents []byte) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(abs), 0750); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	file, err := os.Create(abs)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(contents)
	return err
}

// Read reads bytes from a file
func Read(path string) ([]byte, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, errors.New("trying to read directory as file")
	}

	file, err := os.Open(abs)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}
