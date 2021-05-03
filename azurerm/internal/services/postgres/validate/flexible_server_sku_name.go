package validate

import (
	"fmt"
	"regexp"
)

func FlexibleServerSkuName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if !regexp.MustCompile(`^(B|GP|MO)_((Standard_E(2|4|8|16|32|48|64)s_v3)|(Standard_B(1m|2)s)|(Standard_D(2|4|8|16|32|48|64)s_v3))$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q is not a valid sku name, got %v", k, v))
		return
	}
	return
}
