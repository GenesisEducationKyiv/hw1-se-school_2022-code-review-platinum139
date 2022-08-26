package storage

import (
	"bufio"
	"log"
	"os"
)

type Filestore struct {
	logger   *log.Logger
	filename string
}

func (s *Filestore) Add(record string) error {
	file, err := os.OpenFile(s.filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("Unable to close storage file: %s\n", err)
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

func (s *Filestore) GetAll() ([]string, error) {
	file, err := os.OpenFile(s.filename, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("Unable to close storage file: %s\n", err)
		}
	}()

	var records []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		records = append(records, scanner.Text())
	}

	return records, nil
}

func NewFilestore(logger *log.Logger, filename string) *Filestore {
	return &Filestore{
		logger:   logger,
		filename: filename,
	}
}
