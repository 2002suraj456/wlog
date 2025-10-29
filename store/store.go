package store

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	fileName := fmt.Sprintf("wlog_%d", now.Unix())
	filePath := filepath.Join(s.Options.Directory, fileName)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", filePath, err)
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
		return fmt.Errorf("error writing to file %s: %v", filePath, err)
	}

	return nil
}

func (s *Store) GetFilesList() (fileList []os.DirEntry, err error) {
	files, err := os.ReadDir(s.Options.Directory)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %v", s.Options.Directory, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file)
		}
	}
	return
}

func getFile(filepath string, fileName string) (*os.File, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", fileName, err)
	}
	return file, nil
}

func (s *Store) ReadEntries(fileName string) (lines []string, err error) {
	file, err := getFile(s.Options.Directory, fileName)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing file %s: %v", fileName, err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", fileName, err)
	}
	return
}

// UpdateEntry TODO : testing due
func (s *Store) UpdateEntry(fileName string, id string, newContent []byte) error {
	lines, err := s.ReadEntries(fileName)
	if err != nil {
		return err
	}

	var newlines []byte
	for _, line := range lines {
		parts := strings.SplitN(line, ",", 1)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line format in file %s: %s", fileName, line)
		}
		if parts[0] == id {
			newlines = append(newlines, append([]byte(fmt.Sprintf("%s,", id)), newContent...)...)
		} else {
			newlines = append(newlines, []byte(line)...)
		}
	}
	return s.updateFile(fileName, newlines)
}

func (s *Store) updateFile(fileName string, newlines []byte) error {
	filePath := path.Join(s.Options.Directory, fileName)
	err := os.WriteFile(filePath, newlines, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %v", fileName, err)
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
