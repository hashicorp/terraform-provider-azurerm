package validate

import (
	"fmt"
	"regexp"
)

func SecretRotationInterval(i interface{}, k string) (warnings []string, errors []error) {
	secretRotationInterval := i.(string)

	re := regexp.MustCompile(`^\d+[m|h]$`)
	if re != nil && !re.MatchString(secretRotationInterval) {
		errors = append(errors, fmt.Errorf("%s must be an integer followed by m for minutes or h for hours. Got %q", k, secretRotationInterval))
	}

	return warnings, errors
}