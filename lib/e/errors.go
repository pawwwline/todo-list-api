package e

import "fmt"

func WrapError(message string, err error) error {
	if err != nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}
