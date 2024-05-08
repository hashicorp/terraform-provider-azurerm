// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type RoleManagementPolicyDataSource struct{}

var _ sdk.DataSource = RoleManagementPolicyDataSource{}

type RoleManagementPolicyDataSourceModel struct {
	Scope                   string                                                  `tfschema:"scope"`
	RoleDefinitionId        string                                                  `tfschema:"role_definition_id"`
	Description             *string                                                 `tfschema:"description"`
	DisplayName             *string                                                 `tfschema:"display_name"`
	ActiveAssignmentRules   []RoleManagementPolicyDataSourceActiveAssignmentRules   `tfschema:"active_assignment_rules"`
	EligibleAssignmentRules []RoleManagementPolicyDataSourceEligibleAssignmentRules `tfschema:"eligible_assignment_rules"`
	ActivationRules         []RoleManagementPolicyDataSourceActivationRules         `tfschema:"activation_rules"`
	NotificationRules       []RoleManagementPolicyDataSourceNotificationEvents      `tfschema:"notification_rules"`
}

type RoleManagementPolicyDataSourceActiveAssignmentRules struct {
	ExpirationRequired     bool   `tfschema:"expiration_required"`
	ExpireAfter            string `tfschema:"expire_after"`
	RequireMultiFactorAuth bool   `tfschema:"require_multifactor_authentication"`
	RequireJustification   bool   `tfschema:"require_justification"`
	RequireTicketInfo      bool   `tfschema:"require_ticket_info"`
}

type RoleManagementPolicyDataSourceEligibleAssignmentRules struct {
	ExpirationRequired bool   `tfschema:"expiration_required"`
	ExpireAfter        string `tfschema:"expire_after"`
}

type RoleManagementPolicyDataSourceActivationRules struct {
	MaximumDuration                 string                                        `tfschema:"maximum_duration"`
	RequireApproval                 bool                                          `tfschema:"require_approval"`
	ApprovalStages                  []RoleManagementPolicyDataSourceApprovalStage `tfschema:"approval_stage"`
	RequireConditionalAccessContext string                                        `tfschema:"required_conditional_access_authentication_context"`
	RequireMultiFactorAuth          bool                                          `tfschema:"require_multifactor_authentication"`
	RequireJustification            bool                                          `tfschema:"require_justification"`
	RequireTicketInfo               bool                                          `tfschema:"require_ticket_info"`
}

type RoleManagementPolicyDataSourceApprovalStage struct {
	PrimaryApprovers []RoleManagementPolicyDataSourceApprover `tfschema:"primary_approver"`
}

type RoleManagementPolicyDataSourceApprover struct {
	ID   string `tfschema:"object_id"`
	Type string `tfschema:"type"`
}

type RoleManagementPolicyDataSourceNotificationEvents struct {
	ActiveAssignments   []RoleManagementPolicyDataSourceNotificationRule `tfschema:"active_assignments"`
	EligibleActivations []RoleManagementPolicyDataSourceNotificationRule `tfschema:"eligible_activations"`
	EligibleAssignments []RoleManagementPolicyDataSourceNotificationRule `tfschema:"eligible_assignments"`
}

type RoleManagementPolicyDataSourceNotificationRule struct {
	AdminNotifications    []RoleManagementPolicyDataSourceNotificationSettings `tfschema:"admin_notifications"`
	ApproverNotifications []RoleManagementPolicyDataSourceNotificationSettings `tfschema:"approver_notifications"`
	AssigneeNotifications []RoleManagementPolicyDataSourceNotificationSettings `tfschema:"assignee_notifications"`
}

type RoleManagementPolicyDataSourceNotificationSettings struct {
	NotificationLevel    string   `tfschema:"notification_level"`
	DefaultRecipients    bool     `tfschema:"default_recipients"`
	AdditionalRecipients []string `tfschema:"additional_recipients"`
}

func (r RoleManagementPolicyDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input interface{}, key string) (warnings []string, errors []error) {
		v, ok := input.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", key))
			return
		}

		_, err := parse.RoleManagementPolicyID(v)
		if err != nil {
			errors = append(errors, err)
		}

		return
	}
}

func (r RoleManagementPolicyDataSource) ResourceType() string {
	return "azurerm_role_management_policy"
}

func (r RoleManagementPolicyDataSource) ModelObject() interface{} {
	return &RoleManagementPolicyDataSourceModel{}
}

func (r RoleManagementPolicyDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scope": {
			Description:  "The scope of the role to which this policy will apply",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringStartsWithOneOf("/subscriptions/", "/providers/Microsoft.Management/managementGroups/"),
		},

		"role_definition_id": {
			Description:  "ID of the Azure Role to which this policy is assigned",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile("/providers/Microsoft.Authorization/roleDefinitions/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"), "should be in the format /providers/Microsoft.Authorization/roleDefinitions/00000000-0000-0000-0000-000000000000"),
		},
	}
}

func (r RoleManagementPolicyDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the policy",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"description": {
			Description: "The Description of the policy",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"display_name": {
			Description: "The display name of the policy",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"eligible_assignment_rules": {
			Description: "The rules for eligible assignment of the policy",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expiration_required": {
						Description: "Must the assignment have an expiry date",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"expire_after": {
						Description: "The duration after which assignments expire",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"active_assignment_rules": {
			Description: "The rules for active assignment of the policy",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expiration_required": {
						Description: "Must the assignment have an expiry date",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"expire_after": {
						Description: "The duration after which assignments expire",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"require_multifactor_authentication": {
						Description: "Whether multi-factor authentication is required to make an assignment",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"require_justification": {
						Description: "Whether a justification is required to make an assignment",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"require_ticket_info": {
						Description: "Whether ticket information is required to make an assignment",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},
				},
			},
		},

		"activation_rules": {
			Description: "The activation rules of the policy",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"maximum_duration": {
						Description: "The time after which the an activation can be valid for",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"require_approval": {
						Description: "Whether an approval is required for activation",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"approval_stage": {
						Description: "The approval stages for the activation",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"primary_approver": {
									Description: "The IDs of the users or groups who can approve the activation",
									Type:        pluginsdk.TypeSet,
									Computed:    true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"object_id": {
												Description: "The ID of the object to act as an approver",
												Type:        pluginsdk.TypeString,
												Computed:    true,
											},

											"type": {
												Description: "The type of object acting as an approver",
												Type:        pluginsdk.TypeString,
												Computed:    true,
											},
										},
									},
								},
							},
						},
					},

					"required_conditional_access_authentication_context": {
						Description: "Whether a conditional access context is required during activation",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"require_multifactor_authentication": {
						Description: "Whether multi-factor authentication is required during activation",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"require_justification": {
						Description: "Whether a justification is required during activation",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},

					"require_ticket_info": {
						Description: "Whether ticket information is required during activation",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},
				},
			},
		},

		"notification_rules": {
			Description: "The notification rules of the policy",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"active_assignments": {
						Description: "Notifications about active assignments",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Resource{
							Schema: notificationRuleDataSourceSchema(),
						},
					},
					"eligible_activations": {
						Description: "Notifications about activations of eligible assignments",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Resource{
							Schema: notificationRuleDataSourceSchema(),
						},
					},
					"eligible_assignments": {
						Description: "Notifications about eligible assignments",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Resource{
							Schema: notificationRuleDataSourceSchema(),
						},
					},
				},
			},
		},
	}
}

func (r RoleManagementPolicyDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.RoleManagementPoliciesClient

			var id *parse.RoleManagementPolicyId

			scope := metadata.ResourceData.Get("scope").(string)
			roleDefinitionId := metadata.ResourceData.Get("role_definition_id").(string)

			// We need to find the Assignment to get the Role Definition ID
			assigns, err := metadata.Client.Authorization.RoleManagementPolicyAssignmentsClient.ListForScope(ctx, commonids.NewScopeID(scope))
			if err != nil {
				return fmt.Errorf("failed to list Role Management Policy Assignments for scope %s. %+v", scope, err)
			}

			for _, assignment := range *assigns.Model {
				if *assignment.Properties.RoleDefinitionId == roleDefinitionId {
					id, err = parse.RoleManagementPolicyID(*assignment.Properties.PolicyId)
					if err != nil {
						return err
					}
					break
				}
			}

			resp, err := client.Get(ctx, rolemanagementpolicies.NewScopedRoleManagementPolicyID(id.Scope, id.Name))
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("role management policy %s not found", id.ID())
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := RoleManagementPolicyDataSourceModel{
				Scope:            id.Scope,
				RoleDefinitionId: roleDefinitionId,
			}

			if model := resp.Model; model != nil {
				if prop := model.Properties; prop != nil {
					state.Description = prop.Description
					state.DisplayName = prop.DisplayName

					// Create the rules structure so we can populate them
					if len(state.EligibleAssignmentRules) == 0 {
						state.EligibleAssignmentRules = make([]RoleManagementPolicyDataSourceEligibleAssignmentRules, 1)
					}
					if len(state.ActiveAssignmentRules) == 0 {
						state.ActiveAssignmentRules = make([]RoleManagementPolicyDataSourceActiveAssignmentRules, 1)
					}
					if len(state.ActivationRules) == 0 {
						state.ActivationRules = make([]RoleManagementPolicyDataSourceActivationRules, 1)
					}
					if len(state.NotificationRules) == 0 {
						state.NotificationRules = make([]RoleManagementPolicyDataSourceNotificationEvents, 1)
					}
					if len(state.NotificationRules[0].EligibleActivations) == 0 {
						state.NotificationRules[0].EligibleActivations = make([]RoleManagementPolicyDataSourceNotificationRule, 1)
					}
					if len(state.NotificationRules[0].ActiveAssignments) == 0 {
						state.NotificationRules[0].ActiveAssignments = make([]RoleManagementPolicyDataSourceNotificationRule, 1)
					}
					if len(state.NotificationRules[0].EligibleAssignments) == 0 {
						state.NotificationRules[0].EligibleAssignments = make([]RoleManagementPolicyDataSourceNotificationRule, 1)
					}

					for _, r := range *prop.Rules {
						rule := r.(rolemanagementpolicies.RawRoleManagementPolicyRuleImpl)
						switch rule.Values["id"].(string) {
						case "AuthenticationContext_EndUser_Assignment":
							if claimValue, ok := rule.Values["claimValue"].(string); ok && claimValue != "" {
								state.ActivationRules[0].RequireConditionalAccessContext = claimValue
							}

						case "Approval_EndUser_Assignment":
							if settings, ok := rule.Values["setting"].(map[string]interface{}); ok {
								state.ActivationRules[0].RequireApproval = settings["isApprovalRequired"].(bool)

								if approvalStages, ok := settings["approvalStages"].([]interface{}); ok {
									state.ActivationRules[0].ApprovalStages = make([]RoleManagementPolicyDataSourceApprovalStage, 1)
									approvalStage := approvalStages[0].(map[string]interface{})

									if primaryApprovers, ok := approvalStage["primaryApprovers"].([]interface{}); ok && len(primaryApprovers) > 0 {
										state.ActivationRules[0].ApprovalStages[0].PrimaryApprovers = make([]RoleManagementPolicyDataSourceApprover, len(primaryApprovers))

										for ia, pa := range primaryApprovers {
											approver := pa.(map[string]interface{})
											state.ActivationRules[0].ApprovalStages[0].PrimaryApprovers[ia] = RoleManagementPolicyDataSourceApprover{
												ID:   approver["id"].(string),
												Type: approver["userType"].(string),
											}
										}
									}
								}
							}

						case "Enablement_Admin_Assignment":
							state.ActiveAssignmentRules[0].RequireMultiFactorAuth = false
							state.ActiveAssignmentRules[0].RequireJustification = false

							if enabledRules, ok := rule.Values["enabledRules"].([]interface{}); ok {
								for _, enabledRule := range enabledRules {
									switch enabledRule.(string) {
									case "MultiFactorAuthentication":
										state.ActiveAssignmentRules[0].RequireMultiFactorAuth = true
									case "Justification":
										state.ActiveAssignmentRules[0].RequireJustification = true
									}
								}
							}

						case "Enablement_EndUser_Assignment":
							state.ActivationRules[0].RequireMultiFactorAuth = false
							state.ActivationRules[0].RequireJustification = false
							state.ActivationRules[0].RequireTicketInfo = false

							if enabledRules, ok := rule.Values["enabledRules"].([]interface{}); ok {
								for _, enabledRule := range enabledRules {
									switch enabledRule.(string) {
									case "MultiFactorAuthentication":
										state.ActivationRules[0].RequireMultiFactorAuth = true
									case "Justification":
										state.ActivationRules[0].RequireJustification = true
									case "Ticketing":
										state.ActivationRules[0].RequireTicketInfo = true
									}
								}
							}

						case "Expiration_Admin_Eligibility":
							state.EligibleAssignmentRules[0].ExpirationRequired = rule.Values["isExpirationRequired"].(bool)
							state.EligibleAssignmentRules[0].ExpireAfter = rule.Values["maximumDuration"].(string)

						case "Expiration_Admin_Assignment":
							state.ActiveAssignmentRules[0].ExpirationRequired = rule.Values["isExpirationRequired"].(bool)
							state.ActiveAssignmentRules[0].ExpireAfter = rule.Values["maximumDuration"].(string)

						case "Expiration_EndUser_Assignment":
							state.ActivationRules[0].MaximumDuration = rule.Values["maximumDuration"].(string)

						case "Notification_Admin_Admin_Assignment":
							state.NotificationRules[0].ActiveAssignments[0].AdminNotifications = []RoleManagementPolicyDataSourceNotificationSettings{
								*flattenNotificationDataSourceSettings(pointer.To(rule.Values)),
							}

						case "Notification_Admin_Admin_Eligibility":
							state.NotificationRules[0].EligibleAssignments[0].AdminNotifications = []RoleManagementPolicyDataSourceNotificationSettings{
								*flattenNotificationDataSourceSettings(pointer.To(rule.Values)),
							}

						case "Notification_Admin_EndUser_Assignment":
							state.NotificationRules[0].EligibleActivations[0].AdminNotifications = []RoleManagementPolicyDataSourceNotificationSettings{
								*flattenNotificationDataSourceSettings(pointer.To(rule.Values)),
							}

						case "Notification_Approver_Admin_Assignment":
							state.NotificationRules[0].ActiveAssignments[0].ApproverNotifications = []RoleManagementPolicyDataSourceNotificationSettings{
								*flattenNotificationDataSourceSettings(pointer.To(rule.Values)),
							}

						case "Notification_Approver_Admin_Eligibility":
							state.NotificationRules[0].EligibleAssignments[0].ApproverNotifications = []RoleManagementPolicyDataSourceNotificationSettings{
								*flattenNotificationDataSourceSettings(pointer.To(rule.Values)),
							}

						case "Notification_Approver_EndUser_Assignment":
							state.NotificationRules[0].EligibleActivations[0].ApproverNotifications = []RoleManagementPolicyDataSourceNotificationSettings{
								*flattenNotificationDataSourceSettings(pointer.To(rule.Values)),
							}

						case "Notification_Requestor_Admin_Assignment":
							state.NotificationRules[0].ActiveAssignments[0].AssigneeNotifications = []RoleManagementPolicyDataSourceNotificationSettings{
								*flattenNotificationDataSourceSettings(pointer.To(rule.Values)),
							}

						case "Notification_Requestor_Admin_Eligibility":
							state.NotificationRules[0].EligibleAssignments[0].AssigneeNotifications = []RoleManagementPolicyDataSourceNotificationSettings{
								*flattenNotificationDataSourceSettings(pointer.To(rule.Values)),
							}

						case "Notification_Requestor_EndUser_Assignment":
							state.NotificationRules[0].EligibleActivations[0].AssigneeNotifications = []RoleManagementPolicyDataSourceNotificationSettings{
								*flattenNotificationDataSourceSettings(pointer.To(rule.Values)),
							}
						}
					}
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}

func flattenNotificationDataSourceSettings(rule *map[string]interface{}) *RoleManagementPolicyDataSourceNotificationSettings {
	var notificationLevel string
	var defaultRecipients bool
	var additionalRecipients []string

	if v, ok := (*rule)["notificationLevel"].(string); ok {
		notificationLevel = v
	}
	if v, ok := (*rule)["isDefaultRecipientsEnabled"].(bool); ok {
		defaultRecipients = v
	}
	if v, ok := (*rule)["notificationRecipients"].([]interface{}); ok {
		additionalRecipients = make([]string, len(v))
		for i, r := range v {
			additionalRecipients[i] = r.(string)
		}
	}
	return &RoleManagementPolicyDataSourceNotificationSettings{
		NotificationLevel:    notificationLevel,
		DefaultRecipients:    defaultRecipients,
		AdditionalRecipients: additionalRecipients,
	}
}

func notificationRuleDataSourceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"admin_notifications": {
			Description: "Admin notification settings",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Resource{
				Schema: notificationSettingsDataSourceSchema(),
			},
		},
		"approver_notifications": {
			Description: "Approver notification settings",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Resource{
				Schema: notificationSettingsDataSourceSchema(),
			},
		},
		"assignee_notifications": {
			Description: "Assignee notification settings",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Resource{
				Schema: notificationSettingsDataSourceSchema(),
			},
		},
	}
}

func notificationSettingsDataSourceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"notification_level": {
			Description: "What level of notifications are sent",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"default_recipients": {
			Description: "Whether the default recipients are notified",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},
		"additional_recipients": {
			Description: "The additional recipients to notify",
			Type:        pluginsdk.TypeSet,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
