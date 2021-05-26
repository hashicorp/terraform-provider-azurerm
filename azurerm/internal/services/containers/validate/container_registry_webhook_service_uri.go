package validate

import (
	"fmt"
	"regexp"
)

func ContainerRegistryWebhookServiceUri(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^https?://[^\s]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must start with http:// or https:// and must not contain whitespaces: %q", k, value))
	}

	return warnings, errors
}
