package validate

import (
	"fmt"
	"regexp"
)

func ContainerRegistryWebhookName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9]{5,50}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"alpha numeric characters only are allowed and between 5 and 50 characters in %q: %q", k, value))
	}

	return warnings, errors
}
