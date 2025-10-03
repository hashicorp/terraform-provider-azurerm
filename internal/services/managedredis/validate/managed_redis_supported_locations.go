// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"slices"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
)

func ManagedRedisSupportedLocations(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("%q expected type of to be string", k))
		return
	}

	v = location.Normalize(v)

	// Only limited locations are supported, check the list in https://azure.microsoft.com/explore/global-infrastructure/products-by-region/table
	// under "Redis Enterprise".
	// Note that the API in other locations might give the impression Managed Redis is supported,
	// however subsequent operations (eg: linking geo-replication) then fails / times out.
	locations := []string{
		"australiaeast",
		"brazilsouth",
		"canadacentral",
		"centralindia",
		"centralus",
		"eastasia",
		"eastus",
		"eastus2",
		"germanywestcentral",
		"japaneast",
		"northcentralus",
		"northeurope",
		"southcentralus",
		"southeastasia",
		"swedencentral",
		"uksouth",
		"westeurope",
		"westus",
		"westus2",
		"westus3",
	}

	if !slices.Contains(locations, v) {
		errors = append(errors, fmt.Errorf("location %q does not support Managed Redis, valid locations: %v", v, locations))
		return
	}

	return
}
