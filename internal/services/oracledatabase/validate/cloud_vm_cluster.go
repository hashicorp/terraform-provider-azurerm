package validate

import (
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"
)

func Name(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < 1 || len(v) > 255 {
		errors = append(errors, fmt.Errorf("name must be %d to %d characters", 1, 255))
		return
	}

	firstChar, _ := utf8.DecodeRuneInString(v)
	if !unicode.IsLetter(firstChar) && firstChar != '_' {
		errors = append(errors, fmt.Errorf("name must start with a letter or underscore (_)"))
		return
	}

	re := regexp.MustCompile("--")
	if re.MatchString(v) {
		errors = append(errors, fmt.Errorf("name must not contain any consecutive hyphers (--)"))
		return
	}

	return
}

func CpuCoreCount(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return
	}

	if v < 2 {
		errors = append(errors, fmt.Errorf("cpu_core_count must be at least %v", 2))
		return
	}

	return
}

func DataStorageSizeInTbs(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return
	}

	if v < 2 || v > 192 {
		errors = append(errors, fmt.Errorf("%v must be between %v and %v", k, 2, 192))
		return
	}

	return
}
