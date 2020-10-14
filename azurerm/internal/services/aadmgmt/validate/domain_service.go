package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidateDomainServiceName(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	// The name attribute rules are :
	// 1. The DNS domain name must be less than or equal to 64 characters in length
	// 2. The DNS domain name can only include letters, numbers, periods, and hyphens.
	// 3. The DNS domain name must include at least two segments, and the first segment cannot be all numbers.
	// 4. Each segment of the DNS domain name must start with a letter or number.
	// 5. The prefix of the DNS domain name (for example, 'contoso100' in the DNS domain name 'contoso100.com') must contain 15 or fewer characters.

	if strings.TrimSpace(v) == "" {
		return nil, []error{fmt.Errorf("%q must not be empty", k)}
	}
	if len(v) > 64 {
		return nil, []error{fmt.Errorf("%q must be less than or equal to 64 characters in length", k)}
	}
	if !regexp.MustCompile(`^([a-zA-Z0-9][a-zA-Z0-9\\-]*\.)+[a-zA-Z0-9][a-zA-Z0-9\-]*$`).MatchString(v) {
		return nil, []error{fmt.Errorf("%s can only include letters, numbers, periods, and hyphens, must include at least two segments, Each segment of the DNS domain name must start with a letter or number", k)}
	}

	index := strings.Index(v, ".")
	if index > -1 {
		prefix := v[:index]
		if len(prefix) > 15 || regexp.MustCompile(`^\d+$`).MatchString(prefix) {
			return nil, []error{fmt.Errorf("The prefix of %s (for example, 'contoso100' in the DNS domain name 'contoso100.com') must contain 15 or fewer characters, and can not be all numbers", k)}
		}
	} else {
		return nil, []error{fmt.Errorf("%s must include at least two segments", k)}
	}

	return nil, nil
}
