// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func FlexibleServerSkuName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if !regexp.MustCompile(`^((B_Standard_B((1|2|4|8|12|16|20)ms|2s))|(GP_Standard_D(((2|4|8|16|32|48|64)s_v3)|((2|4|8|16|32|48|64)ds_v4)|((2|4|8|16|32|48|64|96)ds_v5)|((2|4|8|16|32|48|64|96)ads_v5)|(C(2|4|8|16|32|48|64|96)ads_v5)))|(MO_Standard_E((((2|4|8|16|20|32|48|64)s)_v3)|((2|4|6|8|16|20|32|48|64)ds_v4)|((2|4|8|16|20|32|48|64|96)ds_v5)|((2|4|8|16|32|48|64|96)ads_v5)|(C(2|4|8|16|20|32|48|64|96)(ads|as)_v5))))$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q is not a valid sku name, got %v", k, v))
		return
	}
	return
}

func FlexibleServerSkuNameChange(skuOld string, skuNew string) error {
	if len(skuOld) == 0 || len(skuNew) == 0 {
		return nil
	}

	// Migration from a non-confidential to confidential server is currently not supported by Azure.
	// Manual migration process is required. Error and inform consumer.
	skuOldParts := strings.Split(skuOld, "_")
	skuNewParts := strings.Split(skuNew, "_")

	confidentialComputeRegex := "^.C.*$"
	isConfidentialOld := regexp.MustCompile(confidentialComputeRegex).MatchString(skuOldParts[2]) // index 2 is compute part
	isConfidentialNew := regexp.MustCompile(confidentialComputeRegex).MatchString(skuNewParts[2])

	if isConfidentialOld != isConfidentialNew {
		return fmt.Errorf("migration of Postgres flexible server between `non-confidential` and `confidential` compute types requires manual intervention, suggestion is to deploy new server in parallel with old, consult official Azure documentation for migration then remove old server")
	}

	return nil
}
