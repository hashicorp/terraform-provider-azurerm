package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
)

//RedisEnterpriseClusterSkuTypeLocation - validates that the passed sku type is valid in the Redis Enterprist Cluster location or not
func RedisEnterpriseClusterSkuTypeLocation(input string) error {
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
		// NOTE: Yes, all of these locations are currently covered in the cluster location
		//       check but they may diverge in the future so I am checking both
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
