package store

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Option struct {
	Directory string // Directory to store files
}

type Store struct {
	Options Option
}

func NewStore(option Option) *Store {
	return &Store{
		Options: option,
	}
}

func (s *Store) SaveData(content []byte) error {

	if s.Options.Directory == "" {
		return fmt.Errorf("store does not have a directory\nKindly configure the store with a directory")
	}

	if !isDirectory(s.Options.Directory) {
		return fmt.Errorf("store directory %s is not a valid directory", s.Options.Directory)
	}

	now := time.Now()
	fileName := fmt.Sprintf("wlog_%d_%d_%d", now.Year(), now.Month(), now.Day())
	filePath := filepath.Join(s.Options.Directory, fileName)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("error opening file %s: %v", filePath, err))
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing file %s: %v", filePath, err)
		}
	}(file)

	currTime := time.Now().Format(time.RFC3339)
	content = append([]byte(fmt.Sprintf("%s,", currTime)), content...)

	if _, err := file.Write(content); err != nil {
		panic(fmt.Sprintf("error writing to file %s: %v", filePath, err))
	}

	return nil
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
