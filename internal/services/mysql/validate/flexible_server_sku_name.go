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

	// See all available sku names from https://docs.microsoft.com/en-us/azure/mysql/flexible-server/concepts-compute-storage#compute-tiers-size-and-server-types
	if !regexp.MustCompile(`^(B|GP|MO)_((Standard_E(2|4|8|16|32|48|64|80i)ds_v4)|(Standard_E(2|4|8|16|32|48|64|96)ds_v5)|(Standard_B(1|1m|2|2m|4m|8m|12m|16m|20m)s)|(Standard_D(2|4|8|16|32|48|64)ds_v4))$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q is not a valid sku name, got %v", k, v))
		return
	}
	return
}
