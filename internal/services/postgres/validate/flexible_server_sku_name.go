// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FlexibleServerSkuName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^((B_Standard_B((1|2|4|8|12|16|20)ms|2s))|(GP_Standard_D(((2|4|8|16|32|48|64)s_v3)|((2|4|8|16|32|48|64)ds_v4)|((2|4|8|16|32|48|64|96)ds_v5)|((2|4|8|16|32|48|64|96)ds_v6)|((2|4|8|16|32|48|64|96)ads_v5)|(C(2|4|8|16|32|48|64|96)ads_v5)))|(MO_Standard_E((((2|4|8|16|20|32|48|64)s)_v3)|((2|4|6|8|16|20|32|48|64)ds_v4)|((2|4|8|16|20|32|48|64|96)ds_v5)|((2|4|8|16|32|48|64|96)ds_v6)|((2|4|8|16|32|48|64|96)ads_v5)|(C(2|4|8|16|20|32|48|64|96)(ads|as)_v5))))$`), "is not a valid sku name")(i, k)
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
