package azure

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

// NormalizeLocation will be deprecated in the near future, use `location.Normalize()` instead.
func NormalizeLocation(input interface{}) string {
	loc := input.(string)
	return location.Normalize(loc)
}

// SchemaResourceGroupNameDiffSuppress will be deprecated in the near future
// use `commonschema.ResourceGroupName()` instead
func SchemaResourceGroupNameDiffSuppress() *pluginsdk.Schema {
	// @tombuildsstuff: this function should no longer be used, existing resources will need to be worked
	// through (and this switched out for `commonschema.ResourceGroupName()`) once verifying these now pull
	// this value from the Resource ID rather than the API Response.

	return &pluginsdk.Schema{
		Type:             pluginsdk.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: suppress.CaseDifference,
		ValidateFunc:     resourcegroups.ValidateName,
	}
}
