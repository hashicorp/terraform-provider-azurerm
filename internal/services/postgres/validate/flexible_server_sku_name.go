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

	if !regexp.MustCompile(`^((B_Standard_B(1ms|2s))|(GP_Standard_D(((2|4|8|16|32|48|64)s_v3)|((2|4|8|16|32|48|64)ds_v4)))|(MO_Standard_E((((2|4|8|16|20|32|48|64)s)_v3)|((2|4|6|8|16|20|32|48|64)ds_v4))))$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q is not a valid sku name, got %v", k, v))
		return
	}
	return
}
