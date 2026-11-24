// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_group_policy_remediation

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceGroupPolicyRemediationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Policy Remediation.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RemediationName,
		},

		"resource_group_id": {
			Description:  "The Resource Group ID at which the Policy Remediation should be applied.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseResourceGroupIDInsensitively),
		},

		"policy_assignment_id": {
			Description: "The ID of the Policy Assignment that should be remediated.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			// TODO: remove this suppression when github issue https://github.com/Azure/azure-rest-api-specs/issues/8353 is addressed
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validate.PolicyAssignmentID,
		},

		"failure_percentage": {
			Description:  "A number between 0.0 to 1.0 representing the percentage failure threshold. The remediation will fail if the percentage of failed remediation operations (i.e. failed deployments) exceeds this threshold.",
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			ValidateFunc: validation.FloatBetween(0, 1.0),
		},

		"parallel_deployments": {
			Description:  "How many resources to remediate at any given time. Can be used to increase or reduce the pace of the remediation. If not provided, the default parallel deployments value is used.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntPositive,
		},

		"resource_count": {
			Description:  "The max number of resources that can be remediated by the remediation job. If not provided, the default resource count is used.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntPositive,
		},

		"location_filters": {
			Description: "A list of the resource locations that will be remediated.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: location.EnhancedValidate,
			},
		},

		"policy_definition_reference_id": {
			Description: "The unique ID for the policy definition reference within the Policy Set Definition that should be remediated.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// The API does not honour the provided casing, instead returning values in all lowercase
			// https://github.com/Azure/azure-rest-api-specs/issues/37823
			DiffSuppressFunc:      suppress.CaseDifference,
			DiffSuppressOnRefresh: true,
		},

		"resource_discovery_mode": {
			Description:  "The way that resources to remediate are discovered.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(remediations.ResourceDiscoveryModeExistingNonCompliant),
			ValidateFunc: validation.StringInSlice(remediations.PossibleValuesForResourceDiscoveryMode(), false),
		},
	}
}
