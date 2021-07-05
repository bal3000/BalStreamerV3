package errors

import "fmt"

type StatusErr struct {
	StatusCode int
	Message    string
}

func (se StatusErr) Error() string {
	return fmt.Sprintf("%v %s", se.StatusCode, se.Message)
}

func (se StatusErr) Is(target error) bool {
	if status, ok := target.(StatusErr); ok {
		if status.StatusCode == se.StatusCode {
			return true
		}
	}

	return false
}
