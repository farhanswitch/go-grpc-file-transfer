package storage

import (
	"os"
)

type Storage struct {
	dir string
}

func NewStorage(dir string) Storage {
	return Storage{dir}
}
func (s Storage) Store(file *File) error {
	if err := os.WriteFile(s.dir+file.name, file.buffer.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}
