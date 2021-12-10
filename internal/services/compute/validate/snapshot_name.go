package validate

import (
	"fmt"
	"regexp"
)

func SnapshotName(v interface{}, _ string) (warnings []string, errors []error) {
	// a-z, A-Z, 0-9, _ and -. The max name length is 80
	value := v.(string)

	if !regexp.MustCompile("^[A-Za-z0-9_-]+$").MatchString(value) {
		errors = append(errors, fmt.Errorf("Snapshot Names can only contain alphanumeric characters and underscores."))
	}

	length := len(value)
	if length > 80 {
		errors = append(errors, fmt.Errorf("Snapshot Name can be up to 80 characters, currently %d.", length))
	}

	return warnings, errors
}
