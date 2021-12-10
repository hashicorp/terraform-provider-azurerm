package validate

import (
	"fmt"
	"regexp"
)

func TemplateSpecName(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}
	if !regexp.MustCompile(`^[\w\.\-\(\)]{1,64}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s must only contain alpha-numeric characters, parenthesis, underscores, dashes and periods and be between 1 and 64 characters in length", key))
	}
	return
}
