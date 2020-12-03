package validate

import (
	"fmt"
	"net/url"
)

func NetworkConnectionMonitorHttpPath(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) == 0 {
		errors = append(errors, fmt.Errorf("%q cannot be an empty string: %q", k, value))
		return warnings, errors
	}

	path, err := url.ParseRequestURI(value)
	if err != nil {
		errors = append(errors, fmt.Errorf("parsing %q: %q", k, value))
		return warnings, errors
	}

	if path.IsAbs() {
		errors = append(errors, fmt.Errorf("%q only accepts the absolute path: %q", k, value))
		return warnings, errors
	}

	return warnings, errors
}
