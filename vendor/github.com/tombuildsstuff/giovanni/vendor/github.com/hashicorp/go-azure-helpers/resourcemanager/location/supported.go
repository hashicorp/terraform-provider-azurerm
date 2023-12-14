// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package location

import (
	"context"
	"log"
)

// supportedLocations can be (validly) nil - as such this shouldn't be relied on
var supportedLocations *[]string

// CacheSupportedLocations attempts to retrieve the supported locations from the Azure MetaData Service
// and caches them, for used in enhanced validation
func CacheSupportedLocations(ctx context.Context, resourceManagerEndpoint string) {
	locs, err := availableAzureLocations(ctx, resourceManagerEndpoint)
	if err != nil {
		log.Printf("[DEBUG] error retrieving locations: %s. Enhanced validation will be unavailable", err)
		return
	}

	supportedLocations = locs.Locations
}
