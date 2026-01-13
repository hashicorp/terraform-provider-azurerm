package azre001

import (
	"errors"
	"fmt"
)

func validCases() {
	// Valid: errors.New for fixed strings
	_ = errors.New("something went wrong")
	_ = errors.New("invalid input")

	// Valid: fmt.Errorf with placeholders
	value := "test"
	_ = fmt.Errorf("value %s is invalid", value)
	_ = fmt.Errorf("count: %d", 42)
	_ = fmt.Errorf("error: %v", errors.New("nested"))
	_ = fmt.Errorf("wrapped: %w", errors.New("cause"))
}

func invalidCases() {
	// Invalid: fmt.Errorf without placeholders
	_ = fmt.Errorf("something went wrong") // want `AZRE001`
	_ = fmt.Errorf("invalid input")        // want `AZRE001`
	_ = fmt.Errorf("error occurred")       // want `AZRE001`
}
