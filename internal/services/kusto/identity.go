// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusters"
)

func expandTrustedExternalTenants(input []any) *[]clusters.TrustedExternalTenant {
	output := make([]clusters.TrustedExternalTenant, 0)

	for _, v := range input {
		output = append(output, clusters.TrustedExternalTenant{
			Value: pointer.To(v.(string)),
		})
	}

	return &output
}

func flattenTrustedExternalTenants(input *[]clusters.TrustedExternalTenant) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, v := range *input {
		if v.Value == nil {
			continue
		}

		output = append(output, *v.Value)
	}

	return output
}
