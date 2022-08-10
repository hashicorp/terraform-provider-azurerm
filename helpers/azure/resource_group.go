package azure

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

func SchemaResourceGroupName() *pluginsdk.Schema {
	return commonschema.ResourceGroupName()
}

func SchemaResourceGroupNameDiffSuppress() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:             pluginsdk.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: suppress.CaseDifference,
		ValidateFunc:     resourcegroups.ValidateName,
	}
}

// Deprecated: use `commonschema.ResourceGroupNameForDataSource()` instead
func SchemaResourceGroupNameForDataSource() *pluginsdk.Schema {
	return commonschema.ResourceGroupNameForDataSource()
}
