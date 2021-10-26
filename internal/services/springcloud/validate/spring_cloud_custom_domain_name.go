package validate

import (
	"fmt"
	"regexp"
)

func SpringCloudCustomDomainName(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if !regexp.MustCompile(`^([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s should be a valid domain name", v))
	}

	return
}
