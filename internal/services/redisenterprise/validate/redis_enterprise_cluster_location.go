// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

// RedisEnterpriseClusterLocation - validates that the passed interface contains a valid Redis Enterprist Cluster location or not
func RedisEnterpriseClusterLocation(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return warnings, errors
	}

	location := location.Normalize(v)
	validLocations := validRedisEnterpriseClusterLocations()

	for _, str := range validLocations {
		if location == str {
			return warnings, errors
		}
	}

	errors = append(errors, fmt.Errorf("%q does not currently support Redis Enterprise Clusters. Locations which currently support Redis Enterprise Clusters are [%s]", v, azure.QuotedStringSlice(friendlyValidRedisEnterpriseClusterLocations())))
	return warnings, errors
}

func validRedisEnterpriseClusterLocations() []string {
	var validLoc []string

	for _, v := range friendlyValidRedisEnterpriseClusterLocations() {
		validLoc = append(validLoc, location.Normalize(v))
	}

	return validLoc
}

func friendlyValidRedisEnterpriseClusterLocations() []string {
	return []string{
		"Australia East",
		"Australia Southeast",
		"Brazil South",
		"Canada Central",
		"Central India",
		"Central US",
		"Central US EUAP",
		"East Asia",
		"East US",
		"North Central US",
		"North Europe",
		"South Central US",
		"South India",
		"Southeast Asia",
		"UK South",
		"UK West",
		"East US 2",
		"East US 2 EUAP",
		"West Europe",
		"West US",
		"West US 2",
		"West US 3",
	}
}
