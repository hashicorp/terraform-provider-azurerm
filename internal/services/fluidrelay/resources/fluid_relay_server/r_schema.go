// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package fluid_relay_server

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelayservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/fluidrelay/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (s Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Fluid Relay Server.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FluidRelayServerName,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"location":            commonschema.Location(),
		"tags":                commonschema.Tags(),
		"identity":            commonschema.SystemAssignedUserAssignedIdentityOptional(),
		"storage_sku": {
			Description: "SKU of the storage associated with the resource.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C due to the issue linked below
			// todo remove computed when https://github.com/Azure/azure-rest-api-specs/issues/19700 is fixed
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(fluidrelayservers.PossibleValuesForStorageSKU(), false),
		},
		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Description:  "The Key Vault Key ID that will be used to encrypt the Fluid Relay Server.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
					},

					"user_assigned_identity_id": {
						Description:  "The User Assigned Managed Identity ID to be used for accessing the Customer Managed Key for encryption.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},
	}
}

func (s Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"frs_tenant_id": {
			Description: "The Fluid tenantId for this server.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"orderer_endpoints": {
			Description: "An array of the Fluid Relay Orderer endpoints. This will be deprecated in future version of Fluid Relay Server and will always be empty, [more details](https://learn.microsoft.com/azure/azure-fluid-relay/concepts/version-compatibility).",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"storage_endpoints": {
			Description: "An array of storage endpoints for this Fluid Relay Server. This will be deprecated in future version of Fluid Relay Server and will always be empty, [more details](https://learn.microsoft.com/azure/azure-fluid-relay/concepts/version-compatibility).",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"service_endpoints": {
			Description: "An array of service endpoints for this Fluid Relay Server.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"primary_key": {
			Description: "The primary key for this server.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_key": {
			Description: "The secondary key for this server.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
	}
}
