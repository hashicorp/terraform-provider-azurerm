// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func expandTrustedExternalTenants(input []interface{}) *[]clusters.TrustedExternalTenant {
	output := make([]clusters.TrustedExternalTenant, 0)

	for _, v := range input {
		output = append(output, clusters.TrustedExternalTenant{
			Value: utils.String(v.(string)),
		})
	}

	return &output
}

func flattenTrustedExternalTenants(input *[]clusters.TrustedExternalTenant) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		if v.Value == nil {
			continue
		}

		output = append(output, *v.Value)
	}

	return output
}
