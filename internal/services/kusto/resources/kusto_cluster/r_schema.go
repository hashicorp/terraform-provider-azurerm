// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func kustoClusterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ClusterName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForAzureSkuName(),
							false),
					},

					"capacity": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(1, 1000),
					},
				},
			},
		},

		"allowed_fqdns": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"allowed_ip_ranges": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"trusted_external_tenants": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.Any(validation.IsUUID, validation.StringIsEmpty, validation.StringInSlice([]string{"*"}, false)),
			},
		},

		"optimized_auto_scale": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"minimum_instances": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
					"maximum_instances": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
				},
			},
		},

		"uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"data_ingestion_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"language_extension": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForLanguageExtensionName(), false),
					},
					"image": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForLanguageExtensionImageName(), false),
					},
				},
			},
		},

		"public_ip_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(clusters.PublicIPTypeIPvFour),
			ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForPublicIPType(), false),
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"outbound_network_access_restricted": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"double_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"auto_stop_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"disk_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"streaming_ingestion_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"purge_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"tags": commonschema.Tags(),
	}
}
