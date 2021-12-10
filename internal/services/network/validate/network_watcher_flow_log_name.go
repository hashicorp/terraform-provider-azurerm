package validate

import (
	"fmt"
	"regexp"
)

func NetworkWatcherFlowLogName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[^\W_]([\w]{0,79}$|[\w]{0,78}[\w\-.]$)`).MatchString(value) {
		errors = append(errors, fmt.Errorf("the name can be up to 80 characters long. It must begin with a word character, and it must end with a word character or with '_'. The name may contain word characters or '.', '-', '_'. %q: %q", k, value))
	}

	return warnings, errors
}
