package azure

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// Deprecated: use `commonschema.Location()` instead
func SchemaLocation() *pluginsdk.Schema {
	return commonschema.Location()
}

// Deprecated: use `commonschema.ResourceGroupName` instead
func SchemaResourceGroupName() *pluginsdk.Schema {
	return commonschema.ResourceGroupName()
}
