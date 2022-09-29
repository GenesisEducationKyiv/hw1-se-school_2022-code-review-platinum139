package file_storage

import (
	"bufio"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
)

type StorageIntegrationTestSuite struct {
	suite.Suite
}

func (s *StorageIntegrationTestSuite) TestAdd_CreateFileAndAddRecord() {
	// arrange
	logger := log.New(os.Stdout, "", 4)

	filename := "CreateFileAndAddRecordTest.data"
	record := "test_mail@gmail.com"

	err := RemoveFile(filename)
	assert.NoError(s.T(), err)

	store := NewFileStorage(logger, filename)

	// act
	err = store.Add(record)

	// assert
	assert.NoError(s.T(), err)

	exists, err := CheckFileExists(filename)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), true, exists)

	actualRecord, err := ReadRecord(filename)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), record, actualRecord)

	err = RemoveFile(filename)
	assert.NoError(s.T(), err)
}

func (s *StorageIntegrationTestSuite) TestAdd_RecordAlreadyExists() {
	// arrange
	logger := log.New(os.Stdout, "", 4)

	filename := "RecordAlreadyExistsTest.data"
	record := "test_mail@gmail.com"

	err := RemoveFile(filename)
	assert.NoError(s.T(), err)

	err = CreateFile(filename)
	assert.NoError(s.T(), err)

	err = WriteRecord(filename, record)
	assert.NoError(s.T(), err)

	store := NewFileStorage(logger, filename)

	// act
	err = store.Add(record)

	// assert
	assert.Error(s.T(), RecordAlreadyExistsError{}, err)

	err = RemoveFile(filename)
	assert.NoError(s.T(), err)
}

func (s *StorageIntegrationTestSuite) TestGetAll_FileNotExists() {
	// arrange
	logger := log.New(os.Stdout, "", 4)

	filename := "FileNotExistsTest.data"

	err := RemoveFile(filename)
	assert.NoError(s.T(), err)

	store := NewFileStorage(logger, filename)

	// act
	records, err := store.GetAll()

	// assert
	assert.Len(s.T(), records, 0)
	assert.Error(s.T(), err, os.ErrNotExist)
}

func (s *StorageIntegrationTestSuite) TestGetAll_GetOneRecord() {
	// arrange
	logger := log.New(os.Stdout, "", 4)

	filename := "GetOneRecordTest.data"
	record := "test_mail@gmail.com"

	err := RemoveFile(filename)
	assert.NoError(s.T(), err)

	err = CreateFile(filename)
	assert.NoError(s.T(), err)

	err = WriteRecord(filename, record)
	assert.NoError(s.T(), err)

	store := NewFileStorage(logger, filename)

	// act
	actualRecords, err := store.GetAll()

	// assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), actualRecords, 1)
	assert.Equal(s.T(), record, actualRecords[0])

	err = RemoveFile(filename)
	assert.NoError(s.T(), err)
}

func TestStorageIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(StorageIntegrationTestSuite))
}

func CreateFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

func CheckFileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func ReadRecord(filename string) (string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var record string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return record, nil
}

func WriteRecord(filename, record string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(record)
	if err != nil {
		return err
	}
	return nil
}

func RemoveFile(filename string) error {
	err := os.Remove(filename)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
