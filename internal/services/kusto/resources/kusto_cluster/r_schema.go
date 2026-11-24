// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Kusto Cluster to create. Only lowercase Alphanumeric characters allowed, starting with a letter.",
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
						Description: "The name of the SKU.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForAzureSkuName(),
							false),
					},

					"capacity": {
						Description:  "The node count for the cluster. Boundaries depend on the SKU name.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(1, 1000),
					},
				},
			},
		},

		"allowed_fqdns": {
			Description: "List of allowed FQDNs (Fully Qualified Domain Name) for egress from Cluster.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"allowed_ip_ranges": {
			Description: "The list of IPs in the format of CIDR allowed to connect to the cluster.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"trusted_external_tenants": {
			Description: "A list of tenant IDs that are trusted by the cluster. Use `trusted_external_tenants = [\"*\"]` to explicitly allow all other tenants, `trusted_external_tenants = [\"MyTenantOnly\"]` for only your tenant or `trusted_external_tenants = [\"<tenantId1>\", \"<tenantIdx>\"]` to allow specific other tenants.",
			// TODO: Remove this temporary workaround once a global fix for perpetual TypeList order diffs is available.
			// The Azure Data Explorer API returns tenants in its own order, so compare sorted tenant lists.
			DiffSuppressFunc: diffSuppressTrustedExternalTenants,
			Type:             pluginsdk.TypeList,
			Optional:         true,
			Computed:         true, // azignore:AZS007 - pre-existing violation
			ConfigMode:       pluginsdk.SchemaConfigModeAttr,
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
						Description:  "The minimum number of allowed instances.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
					"maximum_instances": {
						Description:  "The maximum number of allowed instances.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
				},
			},
		},

		"uri": {
			Description: "The FQDN of the Azure Kusto Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"data_ingestion_uri": {
			Description: "The Kusto Cluster URI to be used for data ingestion.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"language_extension": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the language extension.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForLanguageExtensionName(), false),
					},
					"image": {
						Description:  "The language extension image.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForLanguageExtensionImageName(), false),
					},
				},
			},
		},

		"public_ip_type": {
			Description:  "What public IP type to create - IPv4 (default), or DualStack (both IPv4 and IPv6).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(clusters.PublicIPTypeIPvFour),
			ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForPublicIPType(), false),
		},

		"public_network_access_enabled": {
			Description: "Is the public network access enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"outbound_network_access_restricted": {
			Description: "Whether to restrict outbound network access.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"double_encryption_enabled": {
			Description: "Is the cluster's double encryption enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"auto_stop_enabled": {
			Description: "Whether the cluster could be automatically stopped (due to lack of data or no activity for many days).",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"disk_encryption_enabled": {
			Description: "Whether the cluster's disks are encrypted.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"streaming_ingestion_enabled": {
			Description: "Whether the streaming ingest is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"purge_enabled": {
			Description: "Whether the purge operations are enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"tags": commonschema.Tags(),
	}
}
