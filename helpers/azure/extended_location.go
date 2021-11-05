package azure

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SchemaExtendedLocation() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Optional:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
}
