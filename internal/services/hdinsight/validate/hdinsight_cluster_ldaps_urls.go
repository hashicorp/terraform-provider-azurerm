// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func HDInsightClusterLdapsUrls(i any, k string) ([]string, []error) {
	return validation.IsURLWithScheme([]string{"ldaps"})(i, k)
}
