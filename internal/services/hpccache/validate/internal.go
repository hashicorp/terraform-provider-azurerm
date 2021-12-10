package validate

import (
	"fmt"
	"strings"
)

func absolutePath(i interface{}, k string) (warnings []string, errs []error) {
	v, ok := i.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if !strings.HasPrefix(v, "/") {
		errs = append(errs, fmt.Errorf(`%s path should start with "/"`, k))
	}
	return warnings, errs
}

func relativePath(i interface{}, k string) (warnings []string, errs []error) {
	v, ok := i.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if strings.HasPrefix(v, "/") {
		errs = append(errs, fmt.Errorf(`%s path should not start with "/"`, k))
	}
	return warnings, errs
}
