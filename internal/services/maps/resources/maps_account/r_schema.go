// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package maps_account

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2023-06-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mapsAccountSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Azure Maps Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Description:  "The SKU of the Azure Maps Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(accounts.PossibleValuesForName(), false),
		},

		"cors": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed_origins": {
						Description: "A list of origins that should be allowed to make cross-origin calls.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						MinItems:    1,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"local_authentication_enabled": {
			Description: "Is local authentication enabled for this Azure Maps Account? When `false`, all authentication to the Azure Maps data-plane REST API is disabled, except Azure AD authentication.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"data_store": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"unique_name": {
						Description:  "The name given to the linked Storage Account.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"storage_account_id": commonschema.ResourceIDReferenceOptional(&commonids.StorageAccountId{}, "Storage Account"),
				},
			},
		},

		"tags": commonschema.Tags(),

		"x_ms_client_id": {
			Description: "A unique identifier for the Maps Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"primary_access_key": {
			Description: "The primary key used to authenticate and authorize access to the Maps REST APIs.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_access_key": {
			Description: "The secondary key used to authenticate and authorize access to the Maps REST APIs.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
	}
}
