package validate

import (
	"fmt"
	"regexp"
)

func SubscriptionName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if len(v) > 64 || v == "" {
		errors = append(errors, fmt.Errorf("Subscription Name must be between 1 and 64 characters in length"))
	}

	if regexp.MustCompile("[<>;|]").MatchString(v) {
		errors = append(errors, fmt.Errorf("Subsciption Name cannot contain the characters `<`, `>`, `;`, or `|`"))
	}

	return
}
