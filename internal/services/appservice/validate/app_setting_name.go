package validate

import (
	"fmt"
	"regexp"
)

func AppSettingName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", k))
		return
	}

	if matched := regexp.MustCompile(`^[0-9a-zA-Z._]+$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, periods and underscores", k))
	}

	return
}
