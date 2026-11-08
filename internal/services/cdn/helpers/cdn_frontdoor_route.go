// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/routes"
)

func FlattenCdnFrontDoorRouteCustomDomainActivatedResourceArray(input *[]routes.ActivatedResourceReference) ([]any, error) {
	results := make([]any, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	// Normalize these values in the configuration file we know they are valid because they were set on the
	// resource... if these are modified in the portal they will all be lowercased...
	for _, customDomain := range *input {
		if customDomain.Id == nil {
			continue
		}
		id, err := afddomains.ParseCustomDomainIDInsensitively(*customDomain.Id)
		if err != nil {
			return nil, err
		}
		results = append(results, id.ID())
	}

	return results, nil
}
