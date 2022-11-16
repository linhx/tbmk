package common

type DuplicateBmkiError struct {
	Msg string
	Id  string
}

func NewDuplicateBmkiError(msg string, id string) *DuplicateBmkiError {
	return &DuplicateBmkiError{Msg: msg, Id: id}
}

func (e *DuplicateBmkiError) Error() string {
	return e.Msg
}
