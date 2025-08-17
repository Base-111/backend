package database

type NotFoundError struct {
	Description string
}

func NewNotFoundError(err error) *NotFoundError {
	return &NotFoundError{Description: err.Error()}
}

func (e *NotFoundError) Error() string {
	return e.Message() + " " + e.Description
}

func (e *NotFoundError) Message() string {
	return "select error: entity not found in table"
}

type RowExistsError struct {
	Description string
}

func NewRowExistsError(err error) *RowExistsError {
	return &RowExistsError{Description: err.Error()}
}

func (e *RowExistsError) Error() string {
	return e.Message() + " " + e.Description
}

func (e *RowExistsError) Message() string {
	return "insert error: row exists in table"
}
