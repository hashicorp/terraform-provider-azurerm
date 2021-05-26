package formatting

import (
	"fmt"
	"time"
)

// ParseAsDateFormat parses the given nilable string as a time.Time using the specified
// format (for example RFC3339)
func ParseAsDateFormat(input *string, format string) (*time.Time, error) {
	if input == nil {
		return nil, nil
	}

	val, err := time.Parse(format, *input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", *input, err)
	}

	return &val, nil
}
