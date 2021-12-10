package validate

import (
	"fmt"
	"regexp"
)

func StorageTargetName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}
	exp := `^[-0-9a-zA-Z_]{1,31}$`
	p := regexp.MustCompile(exp)
	if !p.MatchString(v) {
		// TODO: make this error message less user hostile
		errors = append(errors, fmt.Errorf(`cache target name doesn't comply with regexp: "%s"`, exp))
	}

	return warnings, errors
}
