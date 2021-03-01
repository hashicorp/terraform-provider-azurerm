package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
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
	return []string{
		location.Normalize("Australia East"),
		location.Normalize("Australia Southeast"),
		location.Normalize("Brazil South"),
		location.Normalize("Central US"),
		location.Normalize("Central US EUAP"),
		location.Normalize("East Asia"),
		location.Normalize("East US"),
		location.Normalize("North Europe"),
		location.Normalize("South Central US"),
		location.Normalize("Southeast Asia"),
		location.Normalize("UK South"),
		location.Normalize("UK West"),
	}
}

func friendlyValidRedisEnterpriseClusterLocations() []string {
	return []string{
		"Australia East",
		"Australia Southeast",
		"Brazil South",
		"Central US",
		"Central US EUAP",
		"East Asia",
		"East US",
		"North Europe",
		"South Central US",
		"Southeast Asia",
		"UK South",
		"UK West",
	}
}
