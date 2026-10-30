// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"slices"

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

func flattenTrustedExternalTenants(input *[]clusters.TrustedExternalTenant) []string {
	if input == nil {
		return []string{}
	}

	output := make([]string, 0)

	for _, v := range *input {
		if v.Value == nil {
			continue
		}

		output = append(output, *v.Value)
	}

	// Use the same tenant ordering for refreshes and imports.
	slices.Sort(output)

	return output
}
