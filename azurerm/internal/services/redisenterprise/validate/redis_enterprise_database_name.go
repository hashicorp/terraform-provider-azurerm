package validate

import (
	"fmt"
)

// RedisEnterpriseDatabaseName - validates that the passed interface contains a valid Redis Enterprise Database Name
func RedisEnterpriseDatabaseName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("%q expected type of to be string", k))
		return
	}

	if v != "default" {
		errors = append(errors, fmt.Errorf("%q is currently limited to 'default' only, got %v", k, v))
		return
	}

	return
}
