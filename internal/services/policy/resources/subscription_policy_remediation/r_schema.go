// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription_policy_remediation

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func subscriptionPolicyRemediationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RemediationName,
		},

		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},

		"policy_assignment_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// TODO: remove this suppression when github issue https://github.com/Azure/azure-rest-api-specs/issues/8353 is addressed
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validate.PolicyAssignmentID,
		},

		"failure_percentage": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			ValidateFunc: validation.FloatBetween(0, 1.0),
		},

		"parallel_deployments": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntPositive,
		},

		"resource_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntPositive,
		},

		"location_filters": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: location.EnhancedValidate,
			},
		},

		"policy_definition_reference_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// The API does not honour the provided casing, instead returning values in all lowercase
			// https://github.com/Azure/azure-rest-api-specs/issues/37823
			DiffSuppressFunc:      suppress.CaseDifference,
			DiffSuppressOnRefresh: true,
		},

		"resource_discovery_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(remediations.ResourceDiscoveryModeExistingNonCompliant),
			ValidateFunc: validation.StringInSlice(remediations.PossibleValuesForResourceDiscoveryMode(), false),
		},
	}
}
