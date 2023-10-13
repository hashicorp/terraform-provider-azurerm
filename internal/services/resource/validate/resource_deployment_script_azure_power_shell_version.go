package validate

import (
	"fmt"
	"regexp"
)

func ResourceDeploymentScriptAzurePowerShellVersion(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	var errors []error
	if matched := regexp.MustCompile(`^\d+\.\d+$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q should be in format x.y", k))
	}

	return nil, errors
}
