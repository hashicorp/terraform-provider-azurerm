package validate

import (
	"fmt"
	"regexp"
)

func UrlTemplate(v interface{}, k string) (warnings []string, errors []error) {
	upstreamURL := v.(string)

	if !regexp.MustCompile(`^https?://[^\s]+$`).MatchString(upstreamURL) {
		errors = append(errors, fmt.Errorf(
			"%q must start with http:// or https:// and must not contain whitespaces: %q", k, upstreamURL))
	}

	return warnings, errors
}
