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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	billingValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/billing/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type RoleManagementPolicyResource struct{}

var _ sdk.Resource = RoleManagementPolicyResource{}

type RoleManagementPolicyModel struct {
	Scope                   string                                        `tfschema:"scope"`
	RoleDefinitionId        string                                        `tfschema:"role_definition_id"`
	Name                    string                                        `tfschema:"name"`
	Description             string                                        `tfschema:"description"`
	ActiveAssignmentRules   []RoleManagementPolicyActiveAssignmentRules   `tfschema:"active_assignment_rules"`
	EligibleAssignmentRules []RoleManagementPolicyEligibleAssignmentRules `tfschema:"eligible_assignment_rules"`
	ActivationRules         []RoleManagementPolicyActivationRules         `tfschema:"activation_rules"`
	NotificationRules       []RoleManagementPolicyNotificationEvents      `tfschema:"notification_rules"`
}

type RoleManagementPolicyActiveAssignmentRules struct {
	ExpirationRequired     bool   `tfschema:"expiration_required"`
	ExpireAfter            string `tfschema:"expire_after"`
	RequireMultiFactorAuth bool   `tfschema:"require_multifactor_authentication"`
	RequireJustification   bool   `tfschema:"require_justification"`
	RequireTicketInfo      bool   `tfschema:"require_ticket_info"`
}

type RoleManagementPolicyEligibleAssignmentRules struct {
	ExpirationRequired bool   `tfschema:"expiration_required"`
	ExpireAfter        string `tfschema:"expire_after"`
}

type RoleManagementPolicyActivationRules struct {
	MaximumDuration                 string                              `tfschema:"maximum_duration"`
	RequireApproval                 bool                                `tfschema:"require_approval"`
	ApprovalStages                  []RoleManagementPolicyApprovalStage `tfschema:"approval_stage"`
	RequireConditionalAccessContext string                              `tfschema:"required_conditional_access_authentication_context"`
	RequireMultiFactorAuth          bool                                `tfschema:"require_multifactor_authentication"`
	RequireJustification            bool                                `tfschema:"require_justification"`
	RequireTicketInfo               bool                                `tfschema:"require_ticket_info"`
}

type RoleManagementPolicyApprovalStage struct {
	PrimaryApprovers []RoleManagementPolicyApprover `tfschema:"primary_approver"`
}

type RoleManagementPolicyApprover struct {
	ID   string `tfschema:"object_id"`
	Type string `tfschema:"type"`
}

type RoleManagementPolicyNotificationEvents struct {
	ActiveAssignments   []RoleManagementPolicyNotificationRule `tfschema:"active_assignments"`
	EligibleActivations []RoleManagementPolicyNotificationRule `tfschema:"eligible_activations"`
	EligibleAssignments []RoleManagementPolicyNotificationRule `tfschema:"eligible_assignments"`
}

type RoleManagementPolicyNotificationRule struct {
	AdminNotifications    []RoleManagementPolicyNotificationSettings `tfschema:"admin_notifications"`
	ApproverNotifications []RoleManagementPolicyNotificationSettings `tfschema:"approver_notifications"`
	AssigneeNotifications []RoleManagementPolicyNotificationSettings `tfschema:"assignee_notifications"`
}

type RoleManagementPolicyNotificationSettings struct {
	NotificationLevel    string   `tfschema:"notification_level"`
	DefaultRecipients    bool     `tfschema:"default_recipients"`
	AdditionalRecipients []string `tfschema:"additional_recipients"`
}

func (r RoleManagementPolicyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return parse.ValidateRoleManagementPolicyId
}

func (r RoleManagementPolicyResource) ResourceType() string {
	return "azurerm_role_management_policy"
}

func (r RoleManagementPolicyResource) ModelObject() interface{} {
	return &RoleManagementPolicyModel{}
}

func (r RoleManagementPolicyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_definition_id": {
			Description:  "ID of the Azure Role to which this policy is assigned",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty, // TODO: validate this once we have consolidated the existing role_definition_id ID types across this package, and can also support those with no scope at all
		},

		"scope": {
			Description: "The scope of the role to which this policy will apply",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.Any(
				// Elevated access for a global admin is needed to assign roles in this scope:
				// https://docs.microsoft.com/en-us/azure/role-based-access-control/elevate-access-global-admin#azure-cli
				// It seems only user account is allowed to be elevated access.
				validation.StringMatch(regexp.MustCompile("/providers/Microsoft.Subscription.*"), "Subscription scope is invalid"),

				billingValidate.EnrollmentID,
				commonids.ValidateManagementGroupID,
				commonids.ValidateSubscriptionID,
				commonids.ValidateResourceGroupID,
				azure.ValidateResourceID,
			),
		},

		"eligible_assignment_rules": {
			Description: "The rules for eligible assignment of the policy",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expiration_required": {
						Description: "Must the assignment have an expiry date",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true,
					},

					"expire_after": {
						Description:  "The duration after which assignments expire",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.StringInSlice([]string{"P15D", "P30D", "P90D", "P180D", "P365D"}, false),
					},
				},
			},
		},

		"active_assignment_rules": {
			Description: "The rules for active assignment of the policy",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expiration_required": {
						Description: "Must the assignment have an expiry date",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true,
					},

					"expire_after": {
						Description:  "The duration after which assignments expire",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.StringInSlice([]string{"P15D", "P30D", "P90D", "P180D", "P365D"}, false),
					},

					"require_multifactor_authentication": {
						Description: "Whether multi-factor authentication is required to make an assignment",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true,
					},

					"require_justification": {
						Description: "Whether a justification is required to make an assignment",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true,
					},

					"require_ticket_info": {
						Description: "Whether ticket information is required to make an assignment",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true,
					},
				},
			},
		},

		"activation_rules": {
			Description: "The activation rules of the policy",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"maximum_duration": {
						Description: "The time after which the an activation can be valid for",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Computed:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"PT30M", "PT1H", "PT1H30M", "PT2H", "PT2H30M", "PT3H", "PT3H30M", "PT4H", "PT4H30M", "PT5H", "PT5H30M", "PT6H",
							"PT6H30M", "PT7H", "PT7H30M", "PT8H", "PT8H30M", "PT9H", "PT9H30M", "PT10H", "PT10H30M", "PT11H", "PT11H30M", "PT12H",
							"PT12H30M", "PT13H", "PT13H30M", "PT14H", "PT14H30M", "PT15H", "PT15H30M", "PT16H", "PT16H30M", "PT17H", "PT17H30M", "PT18H",
							"PT18H30M", "PT19H", "PT19H30M", "PT20H", "PT20H30M", "PT21H", "PT21H30M", "PT22H", "PT22H30M", "PT23H", "PT23H30M", "P1D",
						}, false),
					},

					"require_approval": {
						Description: "Whether an approval is required for activation",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true,
					},

					"approval_stage": {
						Description: "The approval stages for the activation",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						MaxItems:    1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"primary_approver": {
									Description: "The IDs of the users or groups who can approve the activation",
									Type:        pluginsdk.TypeSet,
									Required:    true,
									MinItems:    1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"object_id": {
												Description:  "The ID of the object to act as an approver",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.IsUUID,
											},

											"type": {
												Description:  "The type of object acting as an approver",
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
						Description:   "Whether a conditional access context is required during activation",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						Computed:      true,
						ConflictsWith: []string{"activation_rules.0.require_multifactor_authentication"},
						ValidateFunc:  validation.StringIsNotEmpty,
					},

					"require_multifactor_authentication": {
						Description:   "Whether multi-factor authentication is required during activation",
						Type:          pluginsdk.TypeBool,
						Optional:      true,
						Computed:      true,
						ConflictsWith: []string{"activation_rules.0.required_conditional_access_authentication_context"},
					},

					"require_justification": {
						Description: "Whether a justification is required during activation",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true,
					},

					"require_ticket_info": {
						Description: "Whether ticket information is required during activation",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true,
					},
				},
			},
		},

		"notification_rules": {
			Description: "The notification rules of the policy",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"active_assignments": {
						Description: "Notifications about active assignments",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true,
						MaxItems:    1,
						Elem: &pluginsdk.Resource{
							Schema: notificationRuleSchema(),
						},
					},

					"eligible_activations": {
						Description: "Notifications about activations of eligible assignments",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true,
						MaxItems:    1,
						Elem: &pluginsdk.Resource{
							Schema: notificationRuleSchema(),
						},
					},

					"eligible_assignments": {
						Description: "Notifications about eligible assignments",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true,
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

func (r RoleManagementPolicyResource) Attributes() map[string]*pluginsdk.Schema {
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
	}
}

func (r RoleManagementPolicyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.RoleManagementPoliciesClient

			var config RoleManagementPolicyModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			// Resources already have a default policy, so retrieve its ID in order that we can update it
			policyId, err := FindRoleManagementPolicyId(ctx, metadata.Client.Authorization.RoleManagementPoliciesClient, config.Scope, config.RoleDefinitionId)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *policyId)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", policyId, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: model was nil", policyId)
			}

			model, err := buildRoleManagementPolicyForUpdate(pointer.To(metadata), existing.Model)
			if err != nil {
				return err
			}

			if _, err = client.Update(ctx, *policyId, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", policyId, err)
			}

			// We are using a custom type parse.RoleManagementPolicyId as the ID type for this resource, because the actual
			// resource ID type (ScopedRoleManagementPolicyId) changes each time the policy is updated, so this allows us
			// to search for the latest policy at Read time.
			id := parse.NewRoleManagementPolicyId(config.RoleDefinitionId, config.Scope)

			metadata.SetID(id)
			return nil
		},
	}
}

func (r RoleManagementPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.RoleManagementPoliciesClient

			id, err := parse.RoleManagementPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			policyId, err := FindRoleManagementPolicyId(ctx, metadata.Client.Authorization.RoleManagementPoliciesClient, id.Scope, id.RoleDefinitionId)
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *policyId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := RoleManagementPolicyModel{
				Scope:            id.Scope,
				RoleDefinitionId: id.RoleDefinitionId,
				Name:             policyId.RoleManagementPolicyName,
			}

			if model := resp.Model; model != nil {
				if prop := model.Properties; prop != nil {
					state.Description = pointer.From(prop.Description)

					// Create the rules structure so we can populate them
					if len(state.EligibleAssignmentRules) == 0 {
						state.EligibleAssignmentRules = make([]RoleManagementPolicyEligibleAssignmentRules, 1)
					}
					if len(state.ActiveAssignmentRules) == 0 {
						state.ActiveAssignmentRules = make([]RoleManagementPolicyActiveAssignmentRules, 1)
					}
					if len(state.ActivationRules) == 0 {
						state.ActivationRules = make([]RoleManagementPolicyActivationRules, 1)
					}
					if len(state.NotificationRules) == 0 {
						state.NotificationRules = make([]RoleManagementPolicyNotificationEvents, 1)
					}
					if len(state.NotificationRules[0].EligibleActivations) == 0 {
						state.NotificationRules[0].EligibleActivations = make([]RoleManagementPolicyNotificationRule, 1)
					}
					if len(state.NotificationRules[0].ActiveAssignments) == 0 {
						state.NotificationRules[0].ActiveAssignments = make([]RoleManagementPolicyNotificationRule, 1)
					}
					if len(state.NotificationRules[0].EligibleAssignments) == 0 {
						state.NotificationRules[0].EligibleAssignments = make([]RoleManagementPolicyNotificationRule, 1)
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
									state.ActivationRules[0].ApprovalStages = make([]RoleManagementPolicyApprovalStage, 1)
									approvalStage := approvalStages[0].(map[string]interface{})

									if primaryApprovers, ok := approvalStage["primaryApprovers"].([]interface{}); ok && len(primaryApprovers) > 0 {
										state.ActivationRules[0].ApprovalStages[0].PrimaryApprovers = make([]RoleManagementPolicyApprover, len(primaryApprovers))

										for ia, pa := range primaryApprovers {
											approver := pa.(map[string]interface{})
											state.ActivationRules[0].ApprovalStages[0].PrimaryApprovers[ia] = RoleManagementPolicyApprover{
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
							state.NotificationRules[0].ActiveAssignments[0].AdminNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(rule.Values),
							}

						case "Notification_Admin_Admin_Eligibility":
							state.NotificationRules[0].EligibleAssignments[0].AdminNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(rule.Values),
							}

						case "Notification_Admin_EndUser_Assignment":
							state.NotificationRules[0].EligibleActivations[0].AdminNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(rule.Values),
							}

						case "Notification_Approver_Admin_Assignment":
							state.NotificationRules[0].ActiveAssignments[0].ApproverNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(rule.Values),
							}

						case "Notification_Approver_Admin_Eligibility":
							state.NotificationRules[0].EligibleAssignments[0].ApproverNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(rule.Values),
							}

						case "Notification_Approver_EndUser_Assignment":
							state.NotificationRules[0].EligibleActivations[0].ApproverNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(rule.Values),
							}

						case "Notification_Requestor_Admin_Assignment":
							state.NotificationRules[0].ActiveAssignments[0].AssigneeNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(rule.Values),
							}

						case "Notification_Requestor_Admin_Eligibility":
							state.NotificationRules[0].EligibleAssignments[0].AssigneeNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(rule.Values),
							}

						case "Notification_Requestor_EndUser_Assignment":
							state.NotificationRules[0].EligibleActivations[0].AssigneeNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(rule.Values),
							}
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r RoleManagementPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.RoleManagementPoliciesClient

			id, err := parse.RoleManagementPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			policyId, err := FindRoleManagementPolicyId(ctx, metadata.Client.Authorization.RoleManagementPoliciesClient, id.Scope, id.RoleDefinitionId)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *policyId)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", policyId, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			model, err := buildRoleManagementPolicyForUpdate(pointer.To(metadata), existing.Model)
			if err != nil {
				return fmt.Errorf("could not build update request, %+v", err)
			}

			if _, err = client.Update(ctx, *policyId, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", policyId, err)
			}

			return nil
		},
	}
}

func (r RoleManagementPolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if _, err := parse.RoleManagementPolicyID(metadata.ResourceData.Id()); err != nil {
				return err
			}

			// Role Management Policies cannot be deleted, so we'll just return without doing anything
			return nil
		},
	}
}

func notificationRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"admin_notifications": {
			Description: "Admin notification settings",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: notificationSettingsSchema(),
			},
		},

		"approver_notifications": {
			Description: "Approver notification settings",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: notificationSettingsSchema(),
			},
		},

		"assignee_notifications": {
			Description: "Assignee notification settings",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: notificationSettingsSchema(),
			},
		},
	}
}

func notificationSettingsSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"notification_level": {
			Description:  "What level of notifications are sent",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"All", "Critical"}, false),
		},

		"default_recipients": {
			Description: "Whether the default recipients are notified",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},

		"additional_recipients": {
			Description: "The additional recipients to notify",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
