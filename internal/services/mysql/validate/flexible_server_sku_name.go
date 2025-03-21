// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	if !regexp.MustCompile(`^(B_(Standard_B(1m|2|2m|4m|8m|12m|16m|20m)s))|(GP_(Standard_D(2|4|8|16|32|48|64)ds_v4)|(Standard_D(2|4|8|16|32|48|64)ads_v5))|(MO_((Standard_E(2|4|8|16|20|32|48|64|80i)ds_v4)|(Standard_E(2|2a|4|4a|8|8a|16|16a|20|20a|32|32a|48|48a|64|64a|96)ds_v5)))$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q is not a valid sku name, got %v", k, v))
		return
	}
	return
}
