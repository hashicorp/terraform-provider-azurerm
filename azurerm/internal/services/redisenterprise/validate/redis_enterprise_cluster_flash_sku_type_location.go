package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
)

// RedisEnterpriseClusterFlashSkuTypeLocation - validates that the passed location supports the flash sku type or not
func RedisEnterpriseClusterFlashSkuTypeLocation(input string) error {
	location := location.Normalize(input)
	invalidLocations := invalidRedisEnterpriseClusterFlashLocations()

	for _, str := range invalidLocations {
		if location == str {
			return fmt.Errorf("%q does not support Redis Enterprise Clusters Flash SKU's. Locations which do not currently support Redis Enterprise Clusters Flash SKU's are [%s]", input, azure.QuotedStringSlice(friendlyInvalidRedisEnterpriseClusterFlashLocations()))
		}
	}

	return nil
}

func invalidRedisEnterpriseClusterFlashLocations() []string {
	return []string{
		location.Normalize("Australia Southeast"),
		location.Normalize("Brazil South"),
		location.Normalize("Central US"),
		location.Normalize("Central US EUAP"),
		location.Normalize("East Asia"),
		location.Normalize("UK West"),
	}
}

func friendlyInvalidRedisEnterpriseClusterFlashLocations() []string {
	return []string{
		"Australia Southeast",
		"Brazil South",
		"Central US",
		"Central US EUAP",
		"East Asia",
		"UK West",
	}
}
