package validate

import (
	"fmt"
	"net/url"
)

func NetworkConnectionMonitorEndpointAddress(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) == 0 {
		errors = append(errors, fmt.Errorf("%q cannot be an empty string: %q", k, value))
		return warnings, errors
	}

	url, err := url.Parse(value)
	if err != nil {
		errors = append(errors, fmt.Errorf("parsing %q: %q", k, value))
		return warnings, errors
	}

	if url.Scheme != "" || url.RawQuery != "" {
		errors = append(errors, fmt.Errorf("%q cannot contain scheme and query parameter: %q", k, value))
		return warnings, errors
	}

	return warnings, errors
}
