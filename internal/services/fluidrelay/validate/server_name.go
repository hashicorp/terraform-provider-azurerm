package validate

import (
	"fmt"
	"regexp"
)

var serverNameReg = regexp.MustCompile(`^[-0-9a-zA-Z]{1,50}$`)

func FluidRelayServerName(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}
	// Name should contain only alphanumeric characters and hyphens, up to 50 characters long.
	if !serverNameReg.MatchString(v) {
		errors = append(errors, fmt.Errorf("Name should contain only alphanumeric characters and hyphens, up to 50 characters long."))
	}
	return
}
