// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_configuration

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func appConfigurationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the App Configuration.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ConfigurationStoreName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"data_plane_proxy_authentication_mode": {
			Description:  "The data plane proxy authentication mode.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(configurationstores.AuthenticationModeLocal),
			ValidateFunc: validation.StringInSlice(configurationstores.PossibleValuesForAuthenticationMode(), false),
		},

		"data_plane_proxy_private_link_delegation_enabled": {
			Description: "Whether data plane proxy Private Link delegation is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"encryption": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"identity_client_id": {
						Description:  "The client ID of the identity which will be used to access Key Vault.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
					"key_vault_key_identifier": {
						Description:  "The URI of the Key Vault Key used to encrypt data.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"local_auth_enabled": {
			Description: "Whether local authentication methods is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"public_network_access": {
			Description:  "The Public Network Access setting of the App Configuration.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      nil,
			ValidateFunc: validation.StringInSlice(configurationstores.PossibleValuesForPublicNetworkAccess(), true),
		},

		"purge_protection_enabled": {
			Description: "Whether Purge Protection is enabled. This field only works for `standard` SKU.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"replica": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MinItems: 1,
			Set:      helpers.ResourceConfigurationStoreReplicaHash,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the replica.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ConfigurationStoreReplicaName,
					},
					"location": commonschema.LocationWithoutForceNew(),
					"id": {
						Description: "The ID of the App Configuration Replica.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"endpoint": {
						Description: "The URL of the App Configuration Replica.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		// `sku` is not enum, https://github.com/Azure/azure-rest-api-specs/issues/23902
		"sku": {
			Description: "The SKU name of the App Configuration.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "free",
			ValidateFunc: validation.StringInSlice([]string{
				"free",
				"developer",
				"standard",
				"premium",
			}, false),
		},

		"soft_delete_retention_days": {
			Description:  "The number of days that items should be retained for once soft-deleted. This field only works for `standard` SKU.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      7,
			ValidateFunc: validation.IntBetween(1, 7),
			DiffSuppressFunc: func(_, old, new string, _ *pluginsdk.ResourceData) bool {
				return old == "0"
			},
		},

		"tags": commonschema.Tags(),

		"endpoint": {
			Description: "The URL of the App Configuration.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"primary_read_key": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description: "The ID of the Access Key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"connection_string": {
						Description: "The Connection String for this Access Key - consisting of the Endpoint, ID, and Secret.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"secret": {
						Description: "The Secret of the Access Key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
		},

		"primary_write_key": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description: "The ID of the Access Key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"connection_string": {
						Description: "The Connection String for this Access Key - consisting of the Endpoint, ID, and Secret.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"secret": {
						Description: "The Secret of the Access Key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
		},

		"secondary_read_key": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description: "The ID of the Access Key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"connection_string": {
						Description: "The Connection String for this Access Key - consisting of the Endpoint, ID, and Secret.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"secret": {
						Description: "The Secret of the Access Key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
		},

		"secondary_write_key": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description: "The ID of the Access Key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"connection_string": {
						Description: "The Connection String for this Access Key - consisting of the Endpoint, ID, and Secret.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"secret": {
						Description: "The Secret of the Access Key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
		},
	}
}
