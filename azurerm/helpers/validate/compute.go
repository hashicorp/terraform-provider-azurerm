package validate

import (
	"fmt"
	"regexp"
)

func SharedImageVersionName(v interface{}, k string) (ws []string, es []error) {
	value := v.(string)

	r, _ := regexp.Compile("^([0-9]\\.[0-9]\\.[0-9])$")
	if !r.MatchString(value) {
		es = append(es, fmt.Errorf("Expected %s to be in the format `1.2.3` but got %q.", k, value))
	}

	return
}
