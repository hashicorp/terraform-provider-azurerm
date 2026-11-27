// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_policy_exemption

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" //nolint:staticcheck
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourcePolicyExemptionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Policy Exemption.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_id": {
			Description:  "The Resource ID where the Policy Exemption should be applied.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"exemption_category": {
			Description:  "The category of this policy exemption.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInEnumSlice(policy.PossibleExemptionCategoryValues(), false),
		},

		"policy_assignment_id": {
			Description:  "The ID of the Policy Assignment to be exempted at the specified Scope.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.PolicyAssignmentID,
		},

		"display_name": {
			Description:  "A friendly display name to use for this Policy Exemption.",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 128),
		},

		"description": {
			Description:  "A description to use for this Policy Exemption.",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 512),
		},

		"policy_definition_reference_ids": {
			Description: "The policy definition reference ID list when the associated policy assignment is an assignment of a Policy Set Definition.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"expires_on": {
			Description:  "The expiration date and time in UTC ISO 8601 format of this policy exemption.",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.ISO8601DateTime,
		},

		"metadata": helpers.MetadataSchema(),
	}
}
