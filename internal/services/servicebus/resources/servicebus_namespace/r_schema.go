// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func servicebusNamespaceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Service Bus Namespace resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NamespaceName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"sku": {
			Description:  "Which tier to use. Setting this field to `Premium` will force the creation of a new resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(namespaces.PossibleValuesForSkuName(), false),
		},

		"capacity": {
			Description:  "The capacity. When `sku` is `Premium`, capacity can be `1`, `2`, `4`, `8` or `16`. When `sku` is `Basic` or `Standard`, capacity can be `0` only.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 4, 8, 16}),
		},

		"premium_messaging_partitions": {
			Description:  "The number of messaging partitions. Only valid when `sku` is `Premium` and the minimum number is `1`.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      0,
			ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 4}),
		},

		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Description:  "The ID of the Key Vault Key which should be used to encrypt the data in this Service Bus Namespace.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
					},

					"identity_id": {
						Description:  "The ID of the User Assigned Identity that has access to the key.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},

					"infrastructure_encryption_enabled": {
						Description: "Used to specify whether to enable Infrastructure Encryption (Double Encryption).",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						ForceNew:    true,
					},
				},
			},
		},

		"local_auth_enabled": {
			Description: "Whether SAS authentication is enabled for the Service Bus namespace.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"public_network_access_enabled": {
			Description: "Whether public network access is enabled for the Service Bus Namespace.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"minimum_tls_version": {
			Description: "The minimum supported TLS version for this Service Bus Namespace.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(namespaces.TlsVersionOnePointTwo),
			ValidateFunc: validation.StringInSlice([]string{
				string(namespaces.TlsVersionOnePointTwo),
				string(namespaces.TlsVersionOnePointThree),
			}, false),
		},

		"default_primary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_primary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_secondary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"network_rule_set": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"default_action": {
						Description:  "The default action for the Network Rule Set.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(namespaces.DefaultActionAllow),
						ValidateFunc: validation.StringInSlice(namespaces.PossibleValuesForDefaultAction(), false),
					},

					"public_network_access_enabled": {
						Description: "Whether to allow traffic over public network.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"ip_rules": {
						Description: "One or more IP Addresses, or CIDR Blocks which should be able to access the Service Bus Namespace.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"trusted_services_allowed": {
						Description: "Whether Azure Services that are known and trusted for this resource type are allowed to bypass the firewall configuration. See [Trusted Microsoft Services](https://github.com/MicrosoftDocs/azure-docs/blob/master/articles/service-bus-messaging/includes/service-bus-trusted-services.md).",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"network_rules": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      networkRuleHash,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"subnet_id": {
									Description:  "The Subnet ID which should be able to access this Service Bus Namespace.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateSubnetID,
									// The subnet ID returned from the service will have `resourceGroup/{resourceGroupName}` all in lower cases...
									DiffSuppressFunc: suppress.CaseDifference,
								},
								"ignore_missing_vnet_service_endpoint": {
									Description: "Whether the Service Bus Namespace Network Rule Set ignores missing Virtual Network Service Endpoint options in the Subnet.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     false,
								},
							},
						},
					},
				},
			},
		},

		"endpoint": {
			Description: "The URL to access the Service Bus Namespace.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
