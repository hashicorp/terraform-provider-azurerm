package azure

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func SchemaLocation() *pluginsdk.Schema {
	return location.Schema()
}

func SchemaLocationOptional() *pluginsdk.Schema {
	return location.SchemaOptional()
}

func SchemaLocationForDataSource() *pluginsdk.Schema {
	return location.SchemaComputed()
}

// azure.NormalizeLocation is a function which normalises human-readable region/location
// names (e.g. "West US") to the values used and returned by the Azure API (e.g. "westus").
// In state we track the API internal version as it is easier to go from the human form
// to the canonical form than the other way around.
func NormalizeLocation(input interface{}) string {
	loc := input.(string)
	return location.Normalize(loc)
}
