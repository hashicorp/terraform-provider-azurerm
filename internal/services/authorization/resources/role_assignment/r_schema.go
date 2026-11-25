// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package role_assignment

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/billing/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func roleAssignmentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The UUID for this Role Assignment. If omitted, a UUID will be generated.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"scope": {
			Description: "The scope at which the Role Assignment applies to, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`, `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`, or `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM`, or `/providers/Microsoft.Management/managementGroups/myMG`.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.Any(
				// Elevated access (aka User Access Administrator role) is needed to assign roles in the following scopes:
				// https://learn.microsoft.com/azure/role-based-access-control/elevate-access-global-admin#azure-cli
				validation.StringMatch(regexp.MustCompile("/"), "Root scope (/) is invalid"),
				validation.StringMatch(regexp.MustCompile("/providers/Microsoft.Subscription.*"), "Subscription scope is invalid"),
				validation.StringMatch(regexp.MustCompile("/providers/Microsoft.Capacity"), "Capacity scope is invalid"),
				validation.StringMatch(regexp.MustCompile("/providers/Microsoft.BillingBenefits"), "BillingBenefits scope is invalid"),

				validate.EnrollmentID,
				commonids.ValidateManagementGroupID,
				commonids.ValidateSubscriptionID,
				commonids.ValidateResourceGroupID,
				azure.ValidateResourceID,
			),
		},

		"role_definition_id": {
			Description:      "The Scoped-ID of the Role Definition.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true, // azignore:AZS007 - pre-existing violation
			ForceNew:         true,
			ExactlyOneOf:     []string{"role_definition_id", "role_definition_name"},
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"role_definition_name": {
			Description:      "The name of a built-in Role.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true, // azignore:AZS007 - pre-existing violation
			ForceNew:         true,
			ExactlyOneOf:     []string{"role_definition_name", "role_definition_id"},
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validation.StringIsNotEmpty,
		},

		"principal_id": {
			Description:  "The ID of the Principal (User, Group or Service Principal) to assign the Role Definition to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"principal_type": {
			Description: "The type of the `principal_id`. It is necessary to explicitly set this attribute when creating Role Assignments if the principal creating the assignment is constrained by ABAC rules that filters on the PrincipalType attribute.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"User",
				"Group",
				"ServicePrincipal",
			}, false),
		},

		"skip_service_principal_aad_check": {
			Description: "If the `principal_id` is a newly provisioned `Service Principal` set this value to `true` to skip the `Azure Active Directory` check which may fail due to replication lag. This argument is only valid if the `principal_id` is a `Service Principal` identity.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"delegated_managed_identity_resource_id": {
			Description:  "The delegated Azure Resource ID which contains a Managed Identity.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"description": {
			Description:  "The description for this Role Assignment.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"condition": {
			Description:  "The condition that limits the resources that the role can be assigned to.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"condition_version": {
			Description: "The version of the condition.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				"1.0",
				"2.0",
			}, false),
		},
	}
}
