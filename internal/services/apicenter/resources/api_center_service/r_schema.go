// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_service

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Description:  "The name of the API Center service.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.Name(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Description:  "The SKU name for the API Center service.",
			Optional:     true,
			Default:      "Free",
			ValidateFunc: validation.StringInSlice([]string{"Free", "Standard"}, false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if rd.Id() != "" && rd.HasChange("sku_name") {
				oldSkuName, currentSkuName := rd.GetChange("sku_name")
				if oldSkuName == "Standard" && currentSkuName != "Standard" {
					if err := rd.ForceNew("sku_name"); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
}
