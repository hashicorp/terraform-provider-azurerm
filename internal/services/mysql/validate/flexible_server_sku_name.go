// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FlexibleServerSkuName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^(B_(Standard_B(1|1m|2|2m|4m|8m|12m|16m|20m)s))|(GP_(Standard_D(2|4|8|16|32|48|64)ds_v4)|(Standard_D(2|4|8|16|32|48|64)ads_v5))|(MO_((Standard_E(2|4|8|16|20|32|48|64|80i)ds_v4)|(Standard_E(2|2a|4|4a|8|8a|16|16a|20|20a|32|32a|48|48a|64|64a|96|96a)ds_v5)))$`), "is not a valid sku name")(i, k)
}
