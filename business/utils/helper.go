package utils

import (
	"errors"
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func WrapErrors(errs []error) error {
	errorStr := ""
	for _, e := range errs {
		errorStr = errorStr + "." + e.Error()
	}
	return errors.New(errorStr)
}
