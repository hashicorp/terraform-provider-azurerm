// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusters"
)

func expandTrustedExternalTenants(input []interface{}) *[]clusters.TrustedExternalTenant {
	output := make([]clusters.TrustedExternalTenant, 0)

	for _, v := range input {
		output = append(output, clusters.TrustedExternalTenant{
			Value: pointer.To(v.(string)),
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

// trustedExternalTenantsEqual reports whether two `trusted_external_tenants` lists contain the same
// set of tenant IDs, ignoring both order and duplicates. It is used to suppress the perpetual diff
// caused by the Azure Data Explorer API returning the tenants in a different order than configured.
// The attribute is semantically set-valued, so only the unique values are compared.
func trustedExternalTenantsEqual(a, b []interface{}) bool {
	setA := make(map[string]struct{}, len(a))
	for _, v := range a {
		setA[v.(string)] = struct{}{}
	}

	setB := make(map[string]struct{}, len(b))
	for _, v := range b {
		setB[v.(string)] = struct{}{}
	}

	if len(setA) != len(setB) {
		return false
	}

	for k := range setA {
		if _, ok := setB[k]; !ok {
			return false
		}
	}

	return true
}
