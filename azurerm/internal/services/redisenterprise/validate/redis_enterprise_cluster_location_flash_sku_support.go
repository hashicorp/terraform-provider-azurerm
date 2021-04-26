package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
)

// RedisEnterpriseClusterLocationFlashSkuSupport - validates that the passed location supports the flash sku type or not
func RedisEnterpriseClusterLocationFlashSkuSupport(input string) error {
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
	var validFlash []string

	for _, v := range friendlyInvalidRedisEnterpriseClusterFlashLocations() {
		validFlash = append(validFlash, location.Normalize(v))
	}

	return validFlash
}

func friendlyInvalidRedisEnterpriseClusterFlashLocations() []string {
	return []string{
		"Australia Southeast",
		"Brazil South",
		"Central US",
		"Central US EUAP",
		"East Asia",
		"UK West",
		"East US 2 EUAP",
		"South India",
	}
}
