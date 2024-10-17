package validate

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func AutonomousDatabaseName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	firstChar, _ := utf8.DecodeRuneInString(v)
	if !unicode.IsLetter(firstChar) {
		errors = append(errors, fmt.Errorf("%v must start with a letter", k))
		return
	}

	for _, r := range v {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			errors = append(errors, fmt.Errorf("%v must contain only letters and numbers", k))
			return
		}
	}

	if len(v) > 30 {
		errors = append(errors, fmt.Errorf("%v must be 30 characers max", k))
		return
	}

	return
}
