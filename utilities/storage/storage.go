package storage

import (
	"os"
)

type Storage struct {
	Dir string
}

func NewStorage(Dir string) Storage {
	return Storage{Dir}
}
func (s Storage) Store(file *File) error {
	if err := os.WriteFile(s.Dir+file.name, file.buffer.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}
