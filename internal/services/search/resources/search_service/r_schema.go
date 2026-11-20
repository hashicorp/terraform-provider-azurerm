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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(services.PossibleValuesForSkuName(), false),
		},

		"replica_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 12),
		},

		"partition_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  1,
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
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"authentication_failure_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(services.PossibleValuesForAadAuthFailureMode(), false),
		},

		"hosting_mode": {
			Type:                  pluginsdk.TypeString,
			Optional:              true,
			ForceNew:              true,
			Default:               string(services.HostingModeDefault),
			ValidateFunc:          validation.StringInSlice(services.PossibleValuesForHostingMode(), true),
			DiffSuppressFunc:      suppress.CaseDifference, // Breaking change introduced in https://github.com/Azure/azure-rest-api-specs/pull/37579 that changed the case of the Enum value
			DiffSuppressOnRefresh: true,
		},

		"customer_managed_key_enforcement_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"customer_managed_key_encryption_compliance_status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"query_keys": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"key": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"semantic_search_sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(services.SearchSemanticSearchFree),
				string(services.SearchSemanticSearchStandard),
			}, false),
		},

		"allowed_ips": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.Any(
					validation.IsIPv4Address,
					validation.IsCIDRIPv4,
				),
			},
		},

		"network_rule_bypass_option": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(services.PossibleValuesForSearchBypass(), false),
			Default:      string(services.SearchBypassNone),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}
