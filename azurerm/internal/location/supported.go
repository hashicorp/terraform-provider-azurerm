package location

import (
	"context"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
)

// supportedLocations can be (validly) nil - as such this shouldn't be relied on
var supportedLocations *[]string

// CacheSupportedLocations attempts to retrieve the supported locations from the Azure MetaData Service
// and caches them, for used in enhanced validation
func CacheSupportedLocations(ctx context.Context, endpoint string) {
	locs, err := sdk.AvailableAzureLocations(ctx, endpoint)
	if err != nil {
		log.Printf("[DEBUG] error retrieving locations: %s. Enhanced validation will be unavailable", err)
		return
	}

	supportedLocations = locs.Locations
}
