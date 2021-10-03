package validate

import (
	"fmt"
	"regexp"
)

func ProjectName(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}
	validName := regexp.MustCompile(`^^[a-zA-Z0-9][a-zA-Z0-9\-_.]+$*$`)
	if !validName.MatchString(v) {
		return nil, []error{fmt.Errorf("%q must start with letters/numbers and can contain letters, numbers, underscores, dashes and periods - got %q", k, v)}
	}
	return nil, nil
}
