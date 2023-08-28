package validate

import (
	"fmt"
	"regexp"
)

func ManagedLustreFileSystemName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}
	p := regexp.MustCompile(`^[0-9a-zA-Z][-0-9a-zA-Z_]{0,78}[0-9a-zA-Z]$`)
	if !p.MatchString(v) {
		errors = append(errors, fmt.Errorf("%q can contain alphanumeric characters, hyphens and underscores and start and end with alphanumeric and has to be between 2 and 80 characters", k))
	}

	return warnings, errors
}
