// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package role_management_policy

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/billing/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_definition_id": {
			Description:  "The scoped Role Definition ID of the role for which this policy will apply.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty, // TODO: validate this once we have consolidated the existing role_definition_id ID types across this package, and can also support those with no scope at all
		},

		"scope": {
			Description: "The scope to which this Role Management Policy will apply. Can refer to a Management Group, a subscription, a Resource Group or a resource.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.Any(
				// Elevated access for a global admin is needed to assign roles in this scope:
				// https://docs.microsoft.com/azure/role-based-access-control/elevate-access-global-admin#azure-cli
				// It seems only user account is allowed to be elevated access.
				validation.StringMatch(regexp.MustCompile("/providers/Microsoft.Subscription.*"), "Subscription scope is invalid"),

				validate.EnrollmentID,
				commonids.ValidateManagementGroupID,
				commonids.ValidateSubscriptionID,
				commonids.ValidateResourceGroupID,
				azure.ValidateResourceID,
			),
		},

		"eligible_assignment_rules": {
			Description: "The rules for eligible assignment of the policy.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expiration_required": {
						Description: "Must an assignment have an expiry date? `false` allows permanent assignment.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
					},

					"expire_after": {
						Description:  "The maximum length of time an assignment can be valid, as an ISO8601 duration.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice([]string{"P15D", "P30D", "P90D", "P180D", "P365D"}, false),
					},
				},
			},
		},

		"active_assignment_rules": {
			Description: "The rules for active assignment of the policy.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expiration_required": {
						Description: "Must an assignment have an expiry date? `false` allows permanent assignment.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
					},

					"expire_after": {
						Description:  "The maximum length of time an assignment can be valid, as an ISO8601 duration.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice([]string{"P15D", "P30D", "P90D", "P180D", "P365D"}, false),
					},

					"require_multifactor_authentication": {
						Description: "Is multi-factor authentication required to create new assignments?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
					},

					"require_justification": {
						Description: "Is a justification required to create new assignments?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
					},

					"require_ticket_info": {
						Description: "Is ticket information required to create new assignments?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
					},
				},
			},
		},

		"activation_rules": {
			Description: "The activation rules of the policy.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"maximum_duration": {
						Description: "The maximum length of time an activated role can be valid, in an ISO8601 Duration format (e.g. `PT8H`). Valid range is `PT30M` to `PT23H30M`, in 30 minute increments, or `PT1D`.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice([]string{
							"PT30M", "PT1H", "PT1H30M", "PT2H", "PT2H30M", "PT3H", "PT3H30M", "PT4H", "PT4H30M", "PT5H", "PT5H30M", "PT6H",
							"PT6H30M", "PT7H", "PT7H30M", "PT8H", "PT8H30M", "PT9H", "PT9H30M", "PT10H", "PT10H30M", "PT11H", "PT11H30M", "PT12H",
							"PT12H30M", "PT13H", "PT13H30M", "PT14H", "PT14H30M", "PT15H", "PT15H30M", "PT16H", "PT16H30M", "PT17H", "PT17H30M", "PT18H",
							"PT18H30M", "PT19H", "PT19H30M", "PT20H", "PT20H30M", "PT21H", "PT21H30M", "PT22H", "PT22H30M", "PT23H", "PT23H30M", "P1D",
						}, false),
					},

					"require_approval": {
						Description: "Is approval required for activation? If `true` an `approval_stage` block must be provided.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
					},

					"approval_stage": {
						Description: "The approval stages for the activation.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						// Note: O+C because when `activation_rules` is specified, there will be an empty "approval_stage" populated by the API.
						Computed: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"primary_approver": {
									Description: "The IDs of the users or groups who can approve the activation.",
									Type:        pluginsdk.TypeSet,
									Required:    true,
									MinItems:    1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"object_id": {
												Description:  "The ID of the object which will act as an approver.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.IsUUID,
											},

											"type": {
												Description:  "The type of object acting as an approver.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice([]string{"User", "Group"}, false),
											},
										},
									},
								},
							},
						},
					},

					"required_conditional_access_authentication_context": {
						Description:   "The Entra ID Conditional Access context that must be present for activation.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						Computed:      true, // azignore:AZS007 - pre-existing violation
						ConflictsWith: []string{"activation_rules.0.require_multifactor_authentication"},
						ValidateFunc:  validation.StringIsNotEmpty,
					},

					"require_multifactor_authentication": {
						Description:   "Is multi-factor authentication required to activate the role?",
						Type:          pluginsdk.TypeBool,
						Optional:      true,
						Computed:      true, // azignore:AZS007 - pre-existing violation
						ConflictsWith: []string{"activation_rules.0.required_conditional_access_authentication_context"},
					},

					"require_justification": {
						Description: "Is a justification required during activation of the role?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
					},

					"require_ticket_info": {
						Description: "Is ticket information required during activation of the role?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
					},
				},
			},
		},

		"notification_rules": {
			Description: "The notification rules of the policy.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"active_assignments": {
						Description: "Notifications about active assignments.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						MaxItems:    1,
						Elem: &pluginsdk.Resource{
							Schema: notificationRuleSchema(),
						},
					},

					"eligible_activations": {
						Description: "Notifications about activations of eligible assignments.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						MaxItems:    1,
						Elem: &pluginsdk.Resource{
							Schema: notificationRuleSchema(),
						},
					},

					"eligible_assignments": {
						Description: "Notifications about eligible assignments.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						MaxItems:    1,
						Elem: &pluginsdk.Resource{
							Schema: notificationRuleSchema(),
						},
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of this policy, which is typically a UUID and may change over time.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"description": {
			Description: "The description of this policy.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
