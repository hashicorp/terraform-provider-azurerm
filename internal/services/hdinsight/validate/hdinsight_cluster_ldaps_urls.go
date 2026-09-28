// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func HDInsightClusterLdapsUrls(i interface{}, k string) ([]string, []error) {
	return validation.IsURLWithScheme([]string{"ldaps"})(i, k)
}
