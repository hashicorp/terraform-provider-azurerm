package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// validation
func ServiceBusNamespaceName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$").MatchString(value) {
		return
	}

	if !regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{4,48}[a-zA-Z0-9]$").MatchString(value) {
		errors = append(errors, fmt.Errorf("%s must contain only letters, numbers, and hyphens. The namespace must start with a letter, and it must end with a letter or number and be between 6 and 50 characters long or be a GUID/UUID", k))
	}
	// The name cannot end with “-“, “-sb“ or “-mgmt“
	if strings.HasSuffix(value, "-") || strings.HasSuffix(value, "-sb") || strings.HasSuffix(value, "-mgmt") {
		errors = append(errors, fmt.Errorf("%q cannot end with a hyphen, -sb, or -mgmt", k))
	}

	return warnings, errors
}
