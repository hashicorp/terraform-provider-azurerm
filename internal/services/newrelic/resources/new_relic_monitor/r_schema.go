// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package new_relic_monitor

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[\w\-]{1,32}$`),
				`name must be between 1 and 32 characters in length and may contain only letters, numbers, hyphens and underscores`,
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"plan": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"effective_date": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						ForceNew:         true,
						DiffSuppressFunc: suppress.RFC3339Time,
						ValidateFunc:     validation.IsRFC3339Time,
					},

					// Enum is removed https://github.com/Azure/azure-rest-api-specs/issues/31093
					"billing_cycle": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  "MONTHLY",
						ValidateFunc: validation.StringInSlice([]string{
							"MONTHLY",
							"WEEKLY",
							"YEARLY",
						}, false),
					},

					"plan_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  "newrelic-pay-as-you-go-free-live",
						ValidateFunc: validation.StringInSlice([]string{
							"newrelic-pay-as-you-go-free-live",
						}, false),
					},

					"usage_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Default:      string(monitors.UsageTypePAYG),
						ValidateFunc: validation.StringInSlice(monitors.PossibleValuesForUsageType(), false),
					},
				},
			},
		},

		"user": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"first_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"last_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"phone_number": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"account_creation_source": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(monitors.AccountCreationSourceLIFTR),
			ValidateFunc: validation.StringInSlice(monitors.PossibleValuesForAccountCreationSource(), false),
		},

		"account_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.SystemAssignedIdentityOptionalForceNew(),

		"ingestion_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"organization_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"org_creation_source": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(monitors.OrgCreationSourceLIFTR),
			ValidateFunc: validation.StringInSlice(monitors.PossibleValuesForOrgCreationSource(), false),
		},

		"user_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
