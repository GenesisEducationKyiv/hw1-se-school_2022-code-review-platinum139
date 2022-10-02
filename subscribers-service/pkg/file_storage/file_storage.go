package file_storage

import (
	"bufio"
	"os"
	"subscribers-service/internal/common"
)

const filePermissions = 0600

type FileStorage struct {
	logger   common.Logger
	filename string
}

func (s *FileStorage) Add(record string) error {
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
		if scanner.Text() == record {
			return RecordAlreadyExistsError{}
		}
	}

	_, err = file.WriteString(record + "\n")
	if err != nil {
		return err
	}

	return nil
}

func (s *FileStorage) GetAll() ([]string, error) {
	file, err := os.OpenFile(s.filename, os.O_RDWR, filePermissions)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			s.logger.Errorf("Unable to close storage file: ", err)
		}
	}()

	var records []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		records = append(records, scanner.Text())
	}

	return records, nil
}

func NewFileStorage(logger common.Logger, filename string) *FileStorage {
	return &FileStorage{
		logger:   logger,
		filename: filename,
	}
}
