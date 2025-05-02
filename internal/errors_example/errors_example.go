package errors_example

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
)

func IsServerDbError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == "SERVER_DB_ERROR" && e.Code == 500
}

func ErrorServerDbError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, "SERVER_DB_ERROR", fmt.Sprintf(format, args...))
}

func IsWrongContext(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == "WRONG_CONTEXT" && e.Code == 504
}

func ErrorWrongContext(format string, args ...interface{}) *errors.Error {
	return errors.New(504, "WRONG_CONTEXT", fmt.Sprintf(format, args...))
}
