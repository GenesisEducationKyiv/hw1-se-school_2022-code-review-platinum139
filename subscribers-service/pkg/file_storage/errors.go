package file_storage

type RecordAlreadyExistsError struct{}

func (e RecordAlreadyExistsError) Error() string {
	return "record already exists"
}
