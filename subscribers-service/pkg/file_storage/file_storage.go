package file_storage

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"subscribers-service/internal/common"
)

const filePermissions = 0600

type FileStorage struct {
	logger   common.Logger
	filename string
}

func (s *FileStorage) Add(elements ...string) error {
	file, err := os.OpenFile(s.filename, os.O_CREATE|os.O_RDWR, filePermissions)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			s.logger.Errorf("Unable to close storage file:", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		existingElems := strings.Split(scanner.Text(), " | ")

		if existingElems[1] == elements[1] {
			return RecordAlreadyExistsError{}
		}
	}

	_, err = file.WriteString(strings.Join(elements, " | ") + "\n")
	if err != nil {
		return err
	}

	return nil
}

func (s *FileStorage) GetAll() ([][]string, error) {
	file, err := os.OpenFile(s.filename, os.O_CREATE|os.O_RDWR, filePermissions)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			s.logger.Errorf("Unable to close storage file: ", err)
		}
	}()

	var records [][]string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		elems := strings.Split(scanner.Text(), " | ")
		records = append(records, elems)
	}

	return records, nil
}

func (s *FileStorage) Delete(substring string) error {
	input, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")

	var resultLines []string
	for _, line := range lines {
		if !strings.Contains(line, substring) {
			resultLines = append(resultLines, line)
		}
	}

	output := strings.Join(resultLines, "\n")
	err = ioutil.WriteFile(s.filename, []byte(output), filePermissions)
	if err != nil {
		return err
	}

	return nil
}

func NewFileStorage(logger common.Logger, filename string) *FileStorage {
	return &FileStorage{
		logger:   logger,
		filename: filename,
	}
}
