// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_service

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func searchServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name which should be used for this Search Service.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku": {
			Description:  "The SKU which should be used for this Search Service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(services.PossibleValuesForSkuName(), false),
		},

		"replica_count": {
			Description:  "The number of replicas which should be created for this Search Service. This field cannot be set when using a `free` SKU ([see the Microsoft documentation](https://learn.microsoft.com/azure/search/search-sku-tier)).",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 12),
		},

		"partition_count": {
			Description: "The number of partitions which should be created. This field cannot be set when using a `free` SKU ([see the Microsoft documentation](https://learn.microsoft.com/azure/search/search-sku-tier)).",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Default:     1,
			ValidateFunc: validation.IntInSlice([]int{
				1,
				2,
				3,
				4,
				6,
				12,
			}),
		},

		"local_authentication_enabled": {
			Description: "Whether the Search Service allows authenticating using API Keys.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"authentication_failure_mode": {
			Description:  "The response that the Search Service should return for requests that fail authentication.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(services.PossibleValuesForAadAuthFailureMode(), false),
		},

		"hosting_mode": {
			Description:           "The hosting mode, which allows for High Density partitions (that allow for up to 1000 indexes), should be supported.",
			Type:                  pluginsdk.TypeString,
			Optional:              true,
			ForceNew:              true,
			Default:               string(services.HostingModeDefault),
			ValidateFunc:          validation.StringInSlice(services.PossibleValuesForHostingMode(), true),
			DiffSuppressFunc:      suppress.CaseDifference, // Breaking change introduced in https://github.com/Azure/azure-rest-api-specs/pull/37579 that changed the case of the Enum value
			DiffSuppressOnRefresh: true,
		},

		"customer_managed_key_enforcement_enabled": {
			Description: "Whether the Search Service should enforce that non-customer resources are encrypted.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"customer_managed_key_encryption_compliance_status": {
			Description: "Describes whether the Search Service is compliant or not with respect to having non-customer encrypted resources. If a service has more than one non-customer encrypted resource and `Enforcement` is `enabled` then the service will be marked as `NonCompliant`. If all the resources are customer encrypted, then the service will be marked as `Compliant`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"endpoint": {
			Description: "The endpoint used to connect to this Search Service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"primary_key": {
			Description: "The Primary Key used for Search Service Administration.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_key": {
			Description: "The Secondary Key used for Search Service Administration.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"query_keys": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The name of this Query Key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"key": {
						Description: "The value of this Query Key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
		},

		"public_network_access_enabled": {
			Description: "Whether Public Network Access is allowed for this resource.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"semantic_search_sku": {
			Description: "The Semantic Search SKU which should be used for this Search Service.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(services.SearchSemanticSearchFree),
				string(services.SearchSemanticSearchStandard),
			}, false),
		},

		"allowed_ips": {
			Description: "A list of inbound IPv4 addresses or CIDRs that are allowed to access the Search Service. If the incoming IP request is from an IP address which is not included in the `allowed_ips` it will be blocked by the Search Service's firewall.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.Any(
					validation.IsIPv4Address,
					validation.IsCIDRIPv4,
				),
			},
		},

		"network_rule_bypass_option": {
			Description:  "Whether to allow trusted Azure services to access a network restricted Search Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(services.SearchBypassNone),
			ValidateFunc: validation.StringInSlice(services.PossibleValuesForSearchBypass(), false),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}
