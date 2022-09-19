package utils

import (
	"fmt"
)

func WrappedError(method, action string, err error) error {
	return fmt.Errorf("%s %s error: %w", method, action, err)
}
