package validate

import (
	"fmt"
	"strings"
)

func RemoteFilePath(v interface{}, k string) (warnings []string, errors []error) {
	val := v.(string)

	if !strings.HasPrefix(val, "/") {
		errors = append(errors, fmt.Errorf("%q must start with `/`", k))
	}

	return warnings, errors
}
