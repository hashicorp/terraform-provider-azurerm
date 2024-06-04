// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type RoleManagementPolicyResource struct{}

var _ sdk.Resource = RoleManagementPolicyResource{}

type RoleManagementPolicyModel struct {
	Scope                   string                                        `tfschema:"scope"`
	RoleDefinitionId        string                                        `tfschema:"role_definition_id"`
	Description             *string                                       `tfschema:"description"`
	DisplayName             *string                                       `tfschema:"display_name"`
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
	return func(input interface{}, key string) (warnings []string, errors []error) {
		v, ok := input.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", key))
			return
		}

		_, err := rolemanagementpolicies.ParseScopedRoleManagementPolicyID(v)
		if err != nil {
			errors = append(errors, err)
		}

		return
	}
}

func (r RoleManagementPolicyResource) ResourceType() string {
	return "azurerm_role_management_policy"
}

func (r RoleManagementPolicyResource) ModelObject() interface{} {
	return &RoleManagementPolicyModel{}
}

func (r RoleManagementPolicyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scope": {
			Description: "The scope of the role to which this policy will apply",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.Any(
				commonids.ValidateManagementGroupID,
				commonids.ValidateSubscriptionID,
				commonids.ValidateResourceGroupID,
			),
		},

		"role_definition_id": {
			Description:  "ID of the Azure Role to which this policy is assigned",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateRoleDefinitionResourceId,
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

		"display_name": {
			Description: "The display name of the policy",
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

			id := rolemanagementpolicies.NewScopedRoleManagementPolicyID(config.Scope, "")
			roleDefinitionId := config.RoleDefinitionId

			scopedId, err := getScopedPolicyId(ctx, metadata, &id, roleDefinitionId)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *scopedId)
			if err != nil && response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Role Management Policy for %s (Scope %s)", roleDefinitionId, id.Scope)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			model, err := buildPolicyForUpdate(pointer.To(metadata), existing.Model)
			if err != nil {
				return fmt.Errorf("could not build update request, %+v", err)
			}

			resp, err := client.Update(ctx, *scopedId, *model)
			if err != nil {
				return fmt.Errorf("could not create assignment schedule request, %+v", err)
			}

			updatedId, err := rolemanagementpolicies.ParseScopedRoleManagementPolicyID(*resp.Model.Id)
			if err != nil {
				return err
			}
			metadata.SetID(updatedId)
			return nil
		},
	}
}

func (r RoleManagementPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.RoleManagementPoliciesClient

			stateId, err := rolemanagementpolicies.ParseScopedRoleManagementPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// We need to find the Assignment to get the Role Definition ID
			assigns, err := metadata.Client.Authorization.RoleManagementPolicyAssignmentsClient.ListForScope(ctx, commonids.NewScopeID(stateId.Scope))
			if err != nil {
				return fmt.Errorf("failed to list Role Management Policy Assignments for scope %s. %+v", stateId.Scope, err)
			}

			var roleDefinitionId string
			for _, assignment := range *assigns.Model {
				if *assignment.Properties.PolicyId == stateId.ID() {
					roleDefinitionId = *assignment.Properties.RoleDefinitionId
					break
				}
			}

			scopedId, err := getScopedPolicyId(ctx, metadata, stateId, roleDefinitionId)
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *scopedId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(stateId)
				}

				return fmt.Errorf("retrieving %s: %+v", stateId, err)
			}

			state := RoleManagementPolicyModel{
				Scope:            stateId.Scope,
				RoleDefinitionId: roleDefinitionId,
			}

			if model := resp.Model; model != nil {
				if prop := model.Properties; prop != nil {
					state.Description = prop.Description
					state.DisplayName = prop.DisplayName

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

			stateId, err := rolemanagementpolicies.ParseScopedRoleManagementPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			roleDefinitionId := metadata.ResourceData.Get("role_definition_id").(string)

			scopedId, err := getScopedPolicyId(ctx, metadata, stateId, roleDefinitionId)
			if err != nil {
				return err
			}

			var config RoleManagementPolicyModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			existing, err := client.Get(ctx, *scopedId)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Role Management Policy for role %q (Scope %q)", config.RoleDefinitionId, config.Scope)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", stateId)
			}

			model, err := buildPolicyForUpdate(pointer.To(metadata), existing.Model)
			if err != nil {
				return fmt.Errorf("could not build update request, %+v", err)
			}

			resp, err := client.Update(ctx, *scopedId, *model)
			if err != nil {
				return fmt.Errorf("could not create assignment schedule request, %+v", err)
			}

			updatedId, err := rolemanagementpolicies.ParseScopedRoleManagementPolicyID(*resp.Model.Id)
			if err != nil {
				return err
			}
			metadata.SetID(updatedId)
			return nil
		},
	}
}

func (r RoleManagementPolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := rolemanagementpolicies.ParseScopedRoleManagementPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Policies can't be deleted, so we'll fake it by marking as gone.
			return metadata.MarkAsGone(id)
		},
	}
}

func buildPolicyForUpdate(metadata *sdk.ResourceMetaData, rolePolicy *rolemanagementpolicies.RoleManagementPolicy) (*rolemanagementpolicies.RoleManagementPolicy, error) {
	var model RoleManagementPolicyModel
	if err := metadata.Decode(&model); err != nil {
		return nil, fmt.Errorf("decoding: %+v", err)
	}

	roleId, err := parse.ParseRoleDefinitionResourceId(model.RoleDefinitionId)
	if err != nil {
		return nil, fmt.Errorf("parsing role definition id: %+v", err)
	}

	_, err = commonids.ParseManagementGroupID(model.Scope)
	if err != nil {
		scopeId, err := commonids.ParseSubscriptionID(model.Scope)
		if err == nil {
			if strings.HasPrefix(roleId.Scope, scopeId.ID()) {
				return nil, fmt.Errorf("role definition id must be in the same subscription as the scope")
			}
		}
	} else if roleId.Scope != "" {
		return nil, fmt.Errorf("role definition must be scoped to a management group")
	}

	// Take the slice of rules and convert it to a map with the ID as the key
	existingRules := make(map[string]rolemanagementpolicies.RawRoleManagementPolicyRuleImpl)
	for _, r := range *rolePolicy.Properties.Rules {
		rule := r.(rolemanagementpolicies.RawRoleManagementPolicyRuleImpl)
		existingRules[rule.Values["id"].(string)] = rule
	}
	updatedRules := make([]rolemanagementpolicies.RoleManagementPolicyRule, 0)

	if metadata.ResourceData.HasChange("eligible_assignment_rules") {
		if expirationAdminEligibility, ok := existingRules["Expiration_Admin_Eligibility"]; ok && expirationAdminEligibility.Values != nil {
			var expirationRequired bool
			if expirationRequiredRaw, ok := expirationAdminEligibility.Values["isExpirationRequired"]; ok {
				expirationRequired = expirationRequiredRaw.(bool)
			}

			var maximumDuration string
			if maximumDurationRaw, ok := expirationAdminEligibility.Values["maximumDuration"]; ok {
				maximumDuration = maximumDurationRaw.(string)
			}

			if len(model.EligibleAssignmentRules) == 1 {
				if expirationRequired != model.EligibleAssignmentRules[0].ExpirationRequired {
					expirationRequired = model.EligibleAssignmentRules[0].ExpirationRequired
				}
				if maximumDuration != model.EligibleAssignmentRules[0].ExpireAfter &&
					model.EligibleAssignmentRules[0].ExpireAfter != "" {
					maximumDuration = model.EligibleAssignmentRules[0].ExpireAfter
				}
			}

			var id, ruleType string
			var target map[string]interface{}
			if idRaw, ok := expirationAdminEligibility.Values["id"]; ok {
				id = idRaw.(string)
			}
			if ruleTypeRaw, ok := expirationAdminEligibility.Values["ruleType"]; ok {
				ruleType = ruleTypeRaw.(string)
			}
			if targetRaw, ok := expirationAdminEligibility.Values["target"]; ok {
				target = targetRaw.(map[string]interface{})
			}

			updatedRules = append(updatedRules, map[string]interface{}{
				"id":                   id,
				"ruleType":             ruleType,
				"target":               target,
				"isExpirationRequired": expirationRequired,
				"maximumDuration":      maximumDuration,
			})
		}
	}

	if metadata.ResourceData.HasChange("active_assignment_rules.0.require_multifactor_authentication") ||
		metadata.ResourceData.HasChange("active_assignment_rules.0.require_justification") ||
		metadata.ResourceData.HasChange("active_assignment_rules.0.require_ticket_info") {
		if enablementAdminEligibility, ok := existingRules["Enablement_Admin_Assignment"]; ok && enablementAdminEligibility.Values != nil {
			enabledRules := make([]string, 0)

			if len(model.ActiveAssignmentRules) == 1 {
				if model.ActiveAssignmentRules[0].RequireMultiFactorAuth {
					enabledRules = append(enabledRules, "MultiFactorAuthentication")
				}
				if model.ActiveAssignmentRules[0].RequireJustification {
					enabledRules = append(enabledRules, "Justification")
				}
				if model.ActiveAssignmentRules[0].RequireTicketInfo {
					enabledRules = append(enabledRules, "Ticketing")
				}
			}

			var id, ruleType string
			var target map[string]interface{}
			if idRaw, ok := enablementAdminEligibility.Values["id"]; ok {
				id = idRaw.(string)
			}
			if ruleTypeRaw, ok := enablementAdminEligibility.Values["ruleType"]; ok {
				ruleType = ruleTypeRaw.(string)
			}
			if targetRaw, ok := enablementAdminEligibility.Values["target"]; ok {
				target = targetRaw.(map[string]interface{})
			}

			updatedRules = append(updatedRules, map[string]interface{}{
				"id":           id,
				"ruleType":     ruleType,
				"target":       target,
				"enabledRules": enabledRules,
			})
		}
	}

	if metadata.ResourceData.HasChange("active_assignment_rules.0.expiration_required") ||
		metadata.ResourceData.HasChange("active_assignment_rules.0.expire_after") {
		if expirationAdminAssignment, ok := existingRules["Expiration_Admin_Assignment"]; ok && expirationAdminAssignment.Values != nil {
			var expirationRequired bool
			if expirationRequiredRaw, ok := expirationAdminAssignment.Values["isExpirationRequired"]; ok {
				expirationRequired = expirationRequiredRaw.(bool)
			}

			var maximumDuration string
			if maximumDurationRaw, ok := expirationAdminAssignment.Values["maximumDuration"]; ok {
				maximumDuration = maximumDurationRaw.(string)
			}

			if len(model.ActiveAssignmentRules) == 1 {
				if expirationRequired != model.ActiveAssignmentRules[0].ExpirationRequired {
					expirationRequired = model.ActiveAssignmentRules[0].ExpirationRequired
				}
				if maximumDuration != model.ActiveAssignmentRules[0].ExpireAfter &&
					model.ActiveAssignmentRules[0].ExpireAfter != "" {
					maximumDuration = model.ActiveAssignmentRules[0].ExpireAfter
				}
			}

			var id, ruleType string
			var target map[string]interface{}
			if idRaw, ok := expirationAdminAssignment.Values["id"]; ok {
				id = idRaw.(string)
			}
			if ruleTypeRaw, ok := expirationAdminAssignment.Values["ruleType"]; ok {
				ruleType = ruleTypeRaw.(string)
			}
			if targetRaw, ok := expirationAdminAssignment.Values["target"]; ok {
				target = targetRaw.(map[string]interface{})
			}

			updatedRules = append(updatedRules, map[string]interface{}{
				"id":                   id,
				"ruleType":             ruleType,
				"target":               target,
				"isExpirationRequired": expirationRequired,
				"maximumDuration":      maximumDuration,
			})
		}
	}

	if metadata.ResourceData.HasChange("activation_rules.0.maximum_duration") {
		if expirationEndUserAssignment, ok := existingRules["Expiration_EndUser_Assignment"]; ok && expirationEndUserAssignment.Values != nil {
			var id, ruleType, maximumDuration string
			var target map[string]interface{}
			if idRaw, ok := expirationEndUserAssignment.Values["id"]; ok {
				id = idRaw.(string)
			}
			if ruleTypeRaw, ok := expirationEndUserAssignment.Values["ruleType"]; ok {
				ruleType = ruleTypeRaw.(string)
			}
			if targetRaw, ok := expirationEndUserAssignment.Values["target"]; ok {
				target = targetRaw.(map[string]interface{})
			}
			if len(model.ActivationRules) == 1 {
				maximumDuration = model.ActivationRules[0].MaximumDuration
			}

			updatedRules = append(updatedRules, map[string]interface{}{
				"id":              id,
				"ruleType":        ruleType,
				"target":          target,
				"maximumDuration": maximumDuration,
			})
		}
	}

	if metadata.ResourceData.HasChange("activation_rules.0.require_approval") ||
		metadata.ResourceData.HasChange("activation_rules.0.approval_stage") {
		if approvalEndUserAssignment, ok := existingRules["Approval_EndUser_Assignment"]; ok && approvalEndUserAssignment.Values != nil {
			if len(model.ActivationRules) == 1 {
				if model.ActivationRules[0].RequireApproval && len(model.ActivationRules[0].ApprovalStages) != 1 {
					return nil, fmt.Errorf("require_approval is true, but no approval_stages are provided")
				}
			}

			var approvalReqd bool
			var approvalStages []map[string]interface{}

			if settingsRaw, ok := approvalEndUserAssignment.Values["setting"]; ok {
				settings := settingsRaw.(map[string]interface{})

				if approvalReqdRaw, ok := settings["isApprovalRequired"]; ok {
					approvalReqd = approvalReqdRaw.(bool)
				}

				if len(model.ActivationRules) == 1 {
					if approvalReqd != model.ActivationRules[0].RequireApproval {
						approvalReqd = model.ActivationRules[0].RequireApproval
					}
				}

				if metadata.ResourceData.HasChange("activation_rules.0.approval_stage") {
					if len(model.ActivationRules) == 1 {
						approvalStages = make([]map[string]interface{}, len(model.ActivationRules[0].ApprovalStages))
						for i, stage := range model.ActivationRules[0].ApprovalStages {
							primaryApprovers := make([]map[string]interface{}, len(stage.PrimaryApprovers))
							for ia, approver := range stage.PrimaryApprovers {
								primaryApprovers[ia] = map[string]interface{}{
									"id":       approver.ID,
									"userType": approver.Type,
								}
							}

							approvalStages[i] = map[string]interface{}{
								"PrimaryApprovers": primaryApprovers,
							}
						}
					}
				} else {
					if approvalStagesRaw, ok := settings["approvalStages"]; ok {
						approvalStages = approvalStagesRaw.([]map[string]interface{})
					}
				}
			}

			var id, ruleType string
			var target map[string]interface{}
			if idRaw, ok := approvalEndUserAssignment.Values["id"]; ok {
				id = idRaw.(string)
			}
			if ruleTypeRaw, ok := approvalEndUserAssignment.Values["ruleType"]; ok {
				ruleType = ruleTypeRaw.(string)
			}
			if targetRaw, ok := approvalEndUserAssignment.Values["target"]; ok {
				target = targetRaw.(map[string]interface{})
			}

			updatedRules = append(updatedRules, map[string]interface{}{
				"id":       id,
				"ruleType": ruleType,
				"target":   target,
				"setting": map[string]interface{}{
					"isApprovalRequired": approvalReqd,
					"approvalStages":     approvalStages,
				},
			})
		}
	}

	if metadata.ResourceData.HasChange("activation_rules.0.required_conditional_access_authentication_context") {
		if authEndUserAssignment, ok := existingRules["AuthenticationContext_EndUser_Assignment"]; ok && authEndUserAssignment.Values != nil {
			var claimValue string
			if claimValueRaw, ok := authEndUserAssignment.Values["claimValue"]; ok {
				claimValue = claimValueRaw.(string)
			}

			var isEnabled bool
			if _, set := metadata.ResourceData.GetOk("activation_rules.0.required_conditional_access_authentication_context"); set {
				isEnabled = true
				if len(model.ActivationRules) == 1 {
					claimValue = model.ActivationRules[0].RequireConditionalAccessContext
				}
			} else {
				isEnabled = false
			}

			var id, ruleType string
			var target map[string]interface{}
			if idRaw, ok := authEndUserAssignment.Values["id"]; ok {
				id = idRaw.(string)
			}
			if ruleTypeRaw, ok := authEndUserAssignment.Values["ruleType"]; ok {
				ruleType = ruleTypeRaw.(string)
			}
			if targetRaw, ok := authEndUserAssignment.Values["target"]; ok {
				target = targetRaw.(map[string]interface{})
			}

			updatedRules = append(updatedRules, map[string]interface{}{
				"id":         id,
				"ruleType":   ruleType,
				"target":     target,
				"isEnabled":  isEnabled,
				"claimValue": claimValue,
			})
		}
	}

	if metadata.ResourceData.HasChange("activation_rules.0.require_multifactor_authentication") ||
		metadata.ResourceData.HasChange("activation_rules.0.require_justification") ||
		metadata.ResourceData.HasChange("activation_rules.0.require_ticket_info") {
		if enablementEndUserAssignment, ok := existingRules["Enablement_EndUser_Assignment"]; ok && enablementEndUserAssignment.Values != nil {
			enabledRules := make([]string, 0)
			if len(model.ActivationRules) == 1 {
				if model.ActivationRules[0].RequireMultiFactorAuth {
					enabledRules = append(enabledRules, "MultiFactorAuthentication")
				}
				if model.ActivationRules[0].RequireJustification {
					enabledRules = append(enabledRules, "Justification")
				}
				if model.ActivationRules[0].RequireTicketInfo {
					enabledRules = append(enabledRules, "Ticketing")
				}
			}

			var id, ruleType string
			var target map[string]interface{}
			if idRaw, ok := enablementEndUserAssignment.Values["id"]; ok {
				id = idRaw.(string)
			}
			if ruleTypeRaw, ok := enablementEndUserAssignment.Values["ruleType"]; ok {
				ruleType = ruleTypeRaw.(string)
			}
			if targetRaw, ok := enablementEndUserAssignment.Values["target"]; ok {
				target = targetRaw.(map[string]interface{})
			}

			updatedRules = append(updatedRules, map[string]interface{}{
				"id":           id,
				"ruleType":     ruleType,
				"target":       target,
				"enabledRules": enabledRules,
			})
		}
	}

	if metadata.ResourceData.HasChange("notification_rules.0.eligible_assignments.0.admin_notifications") {
		if notificationAdminAdminEligibility, ok := existingRules["Notification_Admin_Admin_Eligibility"]; ok && notificationAdminAdminEligibility.Values != nil {
			if len(model.NotificationRules) == 1 {
				if len(model.NotificationRules[0].EligibleAssignments) == 1 {
					if len(model.NotificationRules[0].EligibleAssignments[0].AdminNotifications) == 1 {
						updatedRules = append(updatedRules,
							expandNotificationSettings(
								notificationAdminAdminEligibility,
								model.NotificationRules[0].EligibleAssignments[0].AdminNotifications[0],
								metadata.ResourceData.HasChange("notification_rules.0.eligible_assignments.0.admin_notifications.0.additional_recipients"),
							),
						)
					}
				}
			}
		}
	}

	if metadata.ResourceData.HasChange("notification_rules.0.active_assignments.0.admin_notifications") {
		if notificationAdminAdminAssignment, ok := existingRules["Notification_Admin_Admin_Assignment"]; ok && notificationAdminAdminAssignment.Values != nil {
			if len(model.NotificationRules) == 1 {
				if len(model.NotificationRules[0].ActiveAssignments) == 1 {
					if len(model.NotificationRules[0].ActiveAssignments[0].AdminNotifications) == 1 {
						updatedRules = append(updatedRules,
							expandNotificationSettings(
								notificationAdminAdminAssignment,
								model.NotificationRules[0].ActiveAssignments[0].AdminNotifications[0],
								metadata.ResourceData.HasChange("notification_rules.0.active_assignments.0.admin_notifications.0.additional_recipients"),
							),
						)
					}
				}
			}
		}
	}

	if metadata.ResourceData.HasChange("notification_rules.0.eligible_activations.0.admin_notifications") {
		if notificationAdminEndUserAssignment, ok := existingRules["Notification_Admin_EndUser_Assignment"]; ok && notificationAdminEndUserAssignment.Values != nil {
			if len(model.NotificationRules) == 1 {
				if len(model.NotificationRules[0].EligibleActivations) == 1 {
					if len(model.NotificationRules[0].EligibleActivations[0].AdminNotifications) == 1 {
						updatedRules = append(updatedRules,
							expandNotificationSettings(
								notificationAdminEndUserAssignment,
								model.NotificationRules[0].EligibleActivations[0].AdminNotifications[0],
								metadata.ResourceData.HasChange("notification_rules.0.eligible_activations.0.admin_notifications.0.additional_recipients"),
							),
						)
					}
				}
			}
		}
	}

	if metadata.ResourceData.HasChange("notification_rules.0.eligible_assignments.0.approver_notifications") {
		if notificationApproverAdminEligibility, ok := existingRules["Notification_Approver_Admin_Eligibility"]; ok && notificationApproverAdminEligibility.Values != nil {
			if len(model.NotificationRules) == 1 {
				if len(model.NotificationRules[0].EligibleAssignments) == 1 {
					if len(model.NotificationRules[0].EligibleAssignments[0].ApproverNotifications) == 1 {
						updatedRules = append(updatedRules,
							expandNotificationSettings(
								notificationApproverAdminEligibility,
								model.NotificationRules[0].EligibleAssignments[0].ApproverNotifications[0],
								metadata.ResourceData.HasChange("notification_rules.0.eligible_assignments.0.approver_notifications.0.additional_recipients"),
							),
						)
					}
				}
			}
		}
	}

	if metadata.ResourceData.HasChange("notification_rules.0.active_assignments.0.approver_notifications") {
		if notificationApproverAdminAssignment, ok := existingRules["Notification_Approver_Admin_Assignment"]; ok && notificationApproverAdminAssignment.Values != nil {
			if len(model.NotificationRules) == 1 {
				if len(model.NotificationRules[0].ActiveAssignments) == 1 {
					if len(model.NotificationRules[0].ActiveAssignments[0].ApproverNotifications) == 1 {
						updatedRules = append(updatedRules,
							expandNotificationSettings(
								notificationApproverAdminAssignment,
								model.NotificationRules[0].ActiveAssignments[0].ApproverNotifications[0],
								metadata.ResourceData.HasChange("notification_rules.0.active_assignments.0.approver_notifications.0.additional_recipients"),
							),
						)
					}
				}
			}
		}
	}

	if metadata.ResourceData.HasChange("notification_rules.0.eligible_activations.0.approver_notifications") {
		if notificationApproverEndUserAssignment, ok := existingRules["Notification_Approver_EndUser_Assignment"]; ok && notificationApproverEndUserAssignment.Values != nil {
			if len(model.NotificationRules) == 1 {
				if len(model.NotificationRules[0].EligibleActivations) == 1 {
					if len(model.NotificationRules[0].EligibleActivations[0].ApproverNotifications) == 1 {
						updatedRules = append(updatedRules,
							expandNotificationSettings(
								notificationApproverEndUserAssignment,
								model.NotificationRules[0].EligibleActivations[0].ApproverNotifications[0],
								metadata.ResourceData.HasChange("notification_rules.0.eligible_activations.0.approver_notifications.0.additional_recipients"),
							),
						)
					}
				}
			}
		}
	}

	if metadata.ResourceData.HasChange("notification_rules.0.eligible_assignments.0.assignee_notifications") {
		if notificationRequestorAdminEligibility, ok := existingRules["Notification_Requestor_Admin_Eligibility"]; ok && notificationRequestorAdminEligibility.Values != nil {
			if len(model.NotificationRules) == 1 {
				if len(model.NotificationRules[0].EligibleAssignments) == 1 {
					if len(model.NotificationRules[0].EligibleAssignments[0].AssigneeNotifications) == 1 {
						updatedRules = append(updatedRules,
							expandNotificationSettings(
								notificationRequestorAdminEligibility,
								model.NotificationRules[0].EligibleAssignments[0].AssigneeNotifications[0],
								metadata.ResourceData.HasChange("notification_rules.0.eligible_assignments.0.assignee_notifications.0.additional_recipients"),
							),
						)
					}
				}
			}
		}
	}

	if metadata.ResourceData.HasChange("notification_rules.0.active_assignments.0.assignee_notifications") {
		if notificationRequestorAdminAssignment, ok := existingRules["Notification_Requestor_Admin_Assignment"]; ok && notificationRequestorAdminAssignment.Values != nil {
			if len(model.NotificationRules) == 1 {
				if len(model.NotificationRules[0].ActiveAssignments) == 1 {
					if len(model.NotificationRules[0].ActiveAssignments[0].AssigneeNotifications) == 1 {
						updatedRules = append(updatedRules,
							expandNotificationSettings(
								notificationRequestorAdminAssignment,
								model.NotificationRules[0].ActiveAssignments[0].AssigneeNotifications[0],
								metadata.ResourceData.HasChange("notification_rules.0.active_assignments.0.assignee_notifications.0.additional_recipients"),
							),
						)
					}
				}
			}
		}
	}

	if metadata.ResourceData.HasChange("notification_rules.0.eligible_activations.0.assignee_notifications") {
		if notificationRequestorEndUserAssignment, ok := existingRules["Notification_Requestor_EndUser_Assignment"]; ok && notificationRequestorEndUserAssignment.Values != nil {
			if len(model.NotificationRules) == 1 {
				if len(model.NotificationRules[0].EligibleActivations) == 1 {
					if len(model.NotificationRules[0].EligibleActivations[0].AssigneeNotifications) == 1 {
						updatedRules = append(updatedRules,
							expandNotificationSettings(
								notificationRequestorEndUserAssignment,
								model.NotificationRules[0].EligibleActivations[0].AssigneeNotifications[0],
								metadata.ResourceData.HasChange("notification_rules.0.eligible_activations.0.assignee_notifications.0.additional_recipients"),
							),
						)
					}
				}
			}
		}
	}

	return &rolemanagementpolicies.RoleManagementPolicy{
		Id:   rolePolicy.Id,
		Name: rolePolicy.Name,
		Type: rolePolicy.Type,
		Properties: &rolemanagementpolicies.RoleManagementPolicyProperties{
			Rules: pointer.To(updatedRules),
		},
	}, nil
}

func expandNotificationSettings(rule rolemanagementpolicies.RawRoleManagementPolicyRuleImpl, data RoleManagementPolicyNotificationSettings, recipientChange bool) rolemanagementpolicies.RoleManagementPolicyRule {
	var level string
	if levelRaw, ok := rule.Values["notificationLevel"]; ok {
		level = levelRaw.(string)
	}

	var defaultRecipients bool
	if defaultRecipientsRaw, ok := rule.Values["isDefaultRecipientsEnabled"]; ok {
		defaultRecipients = defaultRecipientsRaw.(bool)
	}

	var additionalRecipients []string
	if v, ok := rule.Values["notificationRecipients"]; ok {
		additionalRecipientsRaw := v.([]interface{})
		additionalRecipients = make([]string, len(additionalRecipientsRaw))
		for i, r := range additionalRecipientsRaw {
			additionalRecipients[i] = r.(string)
		}
	}

	if level != data.NotificationLevel {
		level = data.NotificationLevel
	}
	if defaultRecipients != data.DefaultRecipients {
		defaultRecipients = data.DefaultRecipients
	}
	if recipientChange {
		additionalRecipients = data.AdditionalRecipients
	}

	var id, ruleType, recipientType, notificationType string
	var target map[string]interface{}
	if idRaw, ok := rule.Values["id"]; ok {
		id = idRaw.(string)
	}
	if ruleTypeRaw, ok := rule.Values["ruleType"]; ok {
		ruleType = ruleTypeRaw.(string)
	}
	if targetRaw, ok := rule.Values["target"]; ok {
		target = targetRaw.(map[string]interface{})
	}
	if recipientTypeRaw, ok := rule.Values["recipientType"]; ok {
		recipientType = recipientTypeRaw.(string)
	}
	if notificationTypeRaw, ok := rule.Values["notificationType"]; ok {
		notificationType = notificationTypeRaw.(string)
	}

	return map[string]interface{}{
		"id":                         id,
		"ruleType":                   ruleType,
		"target":                     target,
		"recipientType":              recipientType,
		"notificationType":           notificationType,
		"notificationLevel":          level,
		"isDefaultRecipientsEnabled": defaultRecipients,
		"notificationRecipients":     additionalRecipients,
	}
}

func flattenNotificationSettings(rule map[string]interface{}) *RoleManagementPolicyNotificationSettings {
	if rule == nil {
		return nil
	}

	var notificationLevel string
	var defaultRecipients bool
	var additionalRecipients []string

	if v, ok := rule["notificationLevel"].(string); ok {
		notificationLevel = v
	}
	if v, ok := rule["isDefaultRecipientsEnabled"].(bool); ok {
		defaultRecipients = v
	}
	if v, ok := rule["notificationRecipients"].([]interface{}); ok {
		additionalRecipients = make([]string, len(v))
		for i, r := range v {
			additionalRecipients[i] = r.(string)
		}
	}
	return &RoleManagementPolicyNotificationSettings{
		NotificationLevel:    notificationLevel,
		DefaultRecipients:    defaultRecipients,
		AdditionalRecipients: additionalRecipients,
	}
}

func getScopedPolicyId(ctx context.Context, metadata sdk.ResourceMetaData, id *rolemanagementpolicies.ScopedRoleManagementPolicyId, roleDefinitionId string) (*rolemanagementpolicies.ScopedRoleManagementPolicyId, error) {
	if id.RoleManagementPolicyName != "" {
		// Check if the current ID is still valid
		resp, err := metadata.Client.Authorization.RoleManagementPoliciesClient.Get(ctx, *id)
		switch {
		case err != nil && !response.WasNotFound(resp.HttpResponse):
			return nil, fmt.Errorf("failed to get existing Role Management Policy %s. %+v", id.ID(), err)
		case resp.Model != nil:
			return id, nil
		}
	}

	// If it wasn't found (or no resource ID was provided), we need to search for the policy ID
	resp, err := metadata.Client.Authorization.RoleManagementPolicyAssignmentsClient.ListForScope(ctx, commonids.NewScopeID(id.Scope))
	if err != nil {
		return nil, fmt.Errorf("failed to list Role Management Policy Assignments for scope %s. %+v", id.Scope, err)
	}

	for _, assignment := range *resp.Model {
		if *assignment.Properties.RoleDefinitionId == roleDefinitionId {
			scopedId, err := rolemanagementpolicies.ParseScopedRoleManagementPolicyID(*assignment.Properties.PolicyId)
			if err != nil {
				return nil, err
			}
			return scopedId, nil
		}
	}

	return nil, fmt.Errorf("could not find Role Management Policy for scope %s and role definition %s", id.Scope, roleDefinitionId)
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
