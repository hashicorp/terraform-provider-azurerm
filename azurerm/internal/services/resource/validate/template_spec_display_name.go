package validate

import "fmt"

func TemplateSpecDisplayName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) == 0 {
		errors = append(errors, fmt.Errorf("%q cannot be an empty string: %q", k, v))
		return
	}

	if len(v) > 64 {
		errors = append(errors, fmt.Errorf("length should be less than %d", 64))
		return
	}

	return
}
