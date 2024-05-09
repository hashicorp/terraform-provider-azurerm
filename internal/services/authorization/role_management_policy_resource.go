// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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
								*flattenNotificationSettings(pointer.To(rule.Values)),
							}

						case "Notification_Admin_Admin_Eligibility":
							state.NotificationRules[0].EligibleAssignments[0].AdminNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(pointer.To(rule.Values)),
							}

						case "Notification_Admin_EndUser_Assignment":
							state.NotificationRules[0].EligibleActivations[0].AdminNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(pointer.To(rule.Values)),
							}

						case "Notification_Approver_Admin_Assignment":
							state.NotificationRules[0].ActiveAssignments[0].ApproverNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(pointer.To(rule.Values)),
							}

						case "Notification_Approver_Admin_Eligibility":
							state.NotificationRules[0].EligibleAssignments[0].ApproverNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(pointer.To(rule.Values)),
							}

						case "Notification_Approver_EndUser_Assignment":
							state.NotificationRules[0].EligibleActivations[0].ApproverNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(pointer.To(rule.Values)),
							}

						case "Notification_Requestor_Admin_Assignment":
							state.NotificationRules[0].ActiveAssignments[0].AssigneeNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(pointer.To(rule.Values)),
							}

						case "Notification_Requestor_Admin_Eligibility":
							state.NotificationRules[0].EligibleAssignments[0].AssigneeNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(pointer.To(rule.Values)),
							}

						case "Notification_Requestor_EndUser_Assignment":
							state.NotificationRules[0].EligibleActivations[0].AssigneeNotifications = []RoleManagementPolicyNotificationSettings{
								*flattenNotificationSettings(pointer.To(rule.Values)),
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

	if strings.HasPrefix(model.Scope, "/subscriptions/") && !strings.HasPrefix(model.RoleDefinitionId, "/subscriptions/") {
		return nil, fmt.Errorf("role_definition_id must be a scoped id")
	}

	// Take the slice of rules and convert it to a map with the ID as the key
	existingRules := make(map[string]rolemanagementpolicies.RawRoleManagementPolicyRuleImpl)
	for _, r := range *rolePolicy.Properties.Rules {
		rule := r.(rolemanagementpolicies.RawRoleManagementPolicyRuleImpl)
		existingRules[rule.Values["id"].(string)] = rule
	}
	updatedRules := make([]rolemanagementpolicies.RoleManagementPolicyRule, 0)

	if metadata.ResourceData.HasChange("eligible_assignment_rules") {
		expirationRequired := existingRules["Expiration_Admin_Eligibility"].Values["isExpirationRequired"].(bool)
		maximumDuration := existingRules["Expiration_Admin_Eligibility"].Values["maximumDuration"].(string)

		if expirationRequired != model.EligibleAssignmentRules[0].ExpirationRequired {
			expirationRequired = model.EligibleAssignmentRules[0].ExpirationRequired
		}
		if maximumDuration != model.EligibleAssignmentRules[0].ExpireAfter &&
			model.EligibleAssignmentRules[0].ExpireAfter != "" {
			maximumDuration = model.EligibleAssignmentRules[0].ExpireAfter
		}

		updatedRules = append(updatedRules, map[string]interface{}{
			"id":                   existingRules["Expiration_Admin_Eligibility"].Values["id"],
			"ruleType":             existingRules["Expiration_Admin_Eligibility"].Values["ruleType"],
			"target":               existingRules["Expiration_Admin_Eligibility"].Values["target"],
			"isExpirationRequired": expirationRequired,
			"maximumDuration":      maximumDuration,
		})
	}

	if metadata.ResourceData.HasChange("active_assignment_rules.0.require_multifactor_authentication") ||
		metadata.ResourceData.HasChange("active_assignment_rules.0.require_justification") ||
		metadata.ResourceData.HasChange("active_assignment_rules.0.require_ticket_info") {
		enabledRules := make([]string, 0)
		if model.ActiveAssignmentRules[0].RequireMultiFactorAuth {
			enabledRules = append(enabledRules, "MultiFactorAuthentication")
		}
		if model.ActiveAssignmentRules[0].RequireJustification {
			enabledRules = append(enabledRules, "Justification")
		}
		if model.ActiveAssignmentRules[0].RequireTicketInfo {
			enabledRules = append(enabledRules, "Ticketing")
		}

		updatedRules = append(updatedRules, map[string]interface{}{
			"id":           existingRules["Enablement_Admin_Assignment"].Values["id"],
			"ruleType":     existingRules["Enablement_Admin_Assignment"].Values["ruleType"],
			"target":       existingRules["Enablement_Admin_Assignment"].Values["target"],
			"enabledRules": enabledRules,
		})
	}

	if metadata.ResourceData.HasChange("active_assignment_rules.0.expiration_required") ||
		metadata.ResourceData.HasChange("active_assignment_rules.0.expire_after") {
		expirationRequired := existingRules["Expiration_Admin_Assignment"].Values["isExpirationRequired"].(bool)
		maximumDuration := existingRules["Expiration_Admin_Assignment"].Values["maximumDuration"].(string)

		if expirationRequired != model.ActiveAssignmentRules[0].ExpirationRequired {
			expirationRequired = model.ActiveAssignmentRules[0].ExpirationRequired
		}
		if maximumDuration != model.ActiveAssignmentRules[0].ExpireAfter &&
			model.ActiveAssignmentRules[0].ExpireAfter != "" {
			maximumDuration = model.ActiveAssignmentRules[0].ExpireAfter
		}

		updatedRules = append(updatedRules, map[string]interface{}{
			"id":                   existingRules["Expiration_Admin_Assignment"].Values["id"],
			"ruleType":             existingRules["Expiration_Admin_Assignment"].Values["ruleType"],
			"target":               existingRules["Expiration_Admin_Assignment"].Values["target"],
			"isExpirationRequired": expirationRequired,
			"maximumDuration":      maximumDuration,
		})
	}

	if metadata.ResourceData.HasChange("activation_rules.0.maximum_duration") {
		updatedRules = append(updatedRules, map[string]interface{}{
			"id":              existingRules["Expiration_EndUser_Assignment"].Values["id"],
			"ruleType":        existingRules["Expiration_EndUser_Assignment"].Values["ruleType"],
			"target":          existingRules["Expiration_EndUser_Assignment"].Values["target"],
			"maximumDuration": model.ActivationRules[0].MaximumDuration,
		})
	}

	if metadata.ResourceData.HasChange("activation_rules.0.require_approval") ||
		metadata.ResourceData.HasChange("activation_rules.0.approval_stage") {
		if model.ActivationRules[0].RequireApproval && len(model.ActivationRules[0].ApprovalStages) != 1 {
			return nil, fmt.Errorf("require_approval is true, but no approval_stages are provided")
		}

		settings := existingRules["Approval_EndUser_Assignment"].Values["setting"].(map[string]interface{})
		approvalReqd := settings["isApprovalRequired"]
		if approvalReqd != model.ActivationRules[0].RequireApproval {
			approvalReqd = model.ActivationRules[0].RequireApproval
		}

		var approvalStages []map[string]interface{}
		if metadata.ResourceData.HasChange("activation_rules.0.approval_stage") {
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
		} else {
			approvalStages = settings["approvalStages"].([]map[string]interface{})
		}

		updatedRules = append(updatedRules, map[string]interface{}{
			"id":       existingRules["Approval_EndUser_Assignment"].Values["id"],
			"ruleType": existingRules["Approval_EndUser_Assignment"].Values["ruleType"],
			"target":   existingRules["Approval_EndUser_Assignment"].Values["target"],
			"setting": map[string]interface{}{
				"isApprovalRequired": approvalReqd,
				"approvalStages":     approvalStages,
			},
		})
	}

	if metadata.ResourceData.HasChange("activation_rules.0.required_conditional_access_authentication_context") {
		var isEnabled bool
		claimValue := existingRules["AuthenticationContext_EndUser_Assignment"].Values["claimValue"]

		if _, set := metadata.ResourceData.GetOk("activation_rules.0.required_conditional_access_authentication_context"); set {
			isEnabled = true
			claimValue = model.ActivationRules[0].RequireConditionalAccessContext
		} else {
			isEnabled = false
		}

		updatedRules = append(updatedRules, map[string]interface{}{
			"id":         existingRules["AuthenticationContext_EndUser_Assignment"].Values["id"],
			"ruleType":   existingRules["AuthenticationContext_EndUser_Assignment"].Values["ruleType"],
			"target":     existingRules["AuthenticationContext_EndUser_Assignment"].Values["target"],
			"isEnabled":  isEnabled,
			"claimValue": claimValue,
		})
	}

	if metadata.ResourceData.HasChange("activation_rules.0.require_multifactor_authentication") ||
		metadata.ResourceData.HasChange("activation_rules.0.require_justification") ||
		metadata.ResourceData.HasChange("activation_rules.0.require_ticket_info") {
		enabledRules := make([]string, 0)
		if model.ActivationRules[0].RequireMultiFactorAuth {
			enabledRules = append(enabledRules, "MultiFactorAuthentication")
		}
		if model.ActivationRules[0].RequireJustification {
			enabledRules = append(enabledRules, "Justification")
		}
		if model.ActivationRules[0].RequireTicketInfo {
			enabledRules = append(enabledRules, "Ticketing")
		}

		updatedRules = append(updatedRules, map[string]interface{}{
			"id":           existingRules["Enablement_EndUser_Assignment"].Values["id"],
			"ruleType":     existingRules["Enablement_EndUser_Assignment"].Values["ruleType"],
			"target":       existingRules["Enablement_EndUser_Assignment"].Values["target"],
			"enabledRules": enabledRules,
		})
	}

	if metadata.ResourceData.HasChange("notification_rules.0.eligible_assignments.0.admin_notifications") {
		updatedRules = append(updatedRules,
			expandNotificationSettings(
				existingRules["Notification_Admin_Admin_Eligibility"],
				model.NotificationRules[0].EligibleAssignments[0].AdminNotifications[0],
				metadata.ResourceData.HasChange("notification_rules.0.eligible_assignments.0.admin_notifications.0.additional_recipients"),
			),
		)
	}

	if metadata.ResourceData.HasChange("notification_rules.0.active_assignments.0.admin_notifications") {
		updatedRules = append(updatedRules,
			expandNotificationSettings(
				existingRules["Notification_Admin_Admin_Assignment"],
				model.NotificationRules[0].ActiveAssignments[0].AdminNotifications[0],
				metadata.ResourceData.HasChange("notification_rules.0.active_assignments.0.admin_notifications.0.additional_recipients"),
			),
		)
	}

	if metadata.ResourceData.HasChange("notification_rules.0.eligible_activations.0.admin_notifications") {
		updatedRules = append(updatedRules,
			expandNotificationSettings(
				existingRules["Notification_Admin_EndUser_Assignment"],
				model.NotificationRules[0].EligibleActivations[0].AdminNotifications[0],
				metadata.ResourceData.HasChange("notification_rules.0.eligible_activations.0.admin_notifications.0.additional_recipients"),
			),
		)
	}

	if metadata.ResourceData.HasChange("notification_rules.0.eligible_assignments.0.approver_notifications") {
		updatedRules = append(updatedRules,
			expandNotificationSettings(
				existingRules["Notification_Approver_Admin_Eligibility"],
				model.NotificationRules[0].EligibleAssignments[0].ApproverNotifications[0],
				metadata.ResourceData.HasChange("notification_rules.0.eligible_assignments.0.approver_notifications.0.additional_recipients"),
			),
		)
	}

	if metadata.ResourceData.HasChange("notification_rules.0.active_assignments.0.approver_notifications") {
		updatedRules = append(updatedRules,
			expandNotificationSettings(
				existingRules["Notification_Approver_Admin_Assignment"],
				model.NotificationRules[0].ActiveAssignments[0].ApproverNotifications[0],
				metadata.ResourceData.HasChange("notification_rules.0.active_assignments.0.approver_notifications.0.additional_recipients"),
			),
		)
	}

	if metadata.ResourceData.HasChange("notification_rules.0.eligible_activations.0.approver_notifications") {
		updatedRules = append(updatedRules,
			expandNotificationSettings(
				existingRules["Notification_Approver_EndUser_Assignment"],
				model.NotificationRules[0].EligibleActivations[0].ApproverNotifications[0],
				metadata.ResourceData.HasChange("notification_rules.0.eligible_activations.0.approver_notifications.0.additional_recipients"),
			),
		)
	}

	if metadata.ResourceData.HasChange("notification_rules.0.eligible_assignments.0.assignee_notifications") {
		updatedRules = append(updatedRules,
			expandNotificationSettings(
				existingRules["Notification_Requestor_Admin_Eligibility"],
				model.NotificationRules[0].EligibleAssignments[0].AssigneeNotifications[0],
				metadata.ResourceData.HasChange("notification_rules.0.eligible_assignments.0.assignee_notifications.0.additional_recipients"),
			),
		)
	}

	if metadata.ResourceData.HasChange("notification_rules.0.active_assignments.0.assignee_notifications") {
		updatedRules = append(updatedRules,
			expandNotificationSettings(
				existingRules["Notification_Requestor_Admin_Assignment"],
				model.NotificationRules[0].ActiveAssignments[0].AssigneeNotifications[0],
				metadata.ResourceData.HasChange("notification_rules.0.active_assignments.0.assignee_notifications.0.additional_recipients"),
			),
		)
	}

	if metadata.ResourceData.HasChange("notification_rules.0.eligible_activations.0.assignee_notifications") {
		updatedRules = append(updatedRules,
			expandNotificationSettings(
				existingRules["Notification_Requestor_EndUser_Assignment"],
				model.NotificationRules[0].EligibleActivations[0].AssigneeNotifications[0],
				metadata.ResourceData.HasChange("notification_rules.0.eligible_activations.0.assignee_notifications.0.additional_recipients"),
			),
		)
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
	level := rule.Values["notificationLevel"]
	defaultRecipients := rule.Values["isDefaultRecipientsEnabled"]
	additionalRecipients := rule.Values["notificationRecipients"]

	if level != data.NotificationLevel {
		level = data.NotificationLevel
	}
	if defaultRecipients != data.DefaultRecipients {
		defaultRecipients = data.DefaultRecipients
	}
	if recipientChange {
		additionalRecipients = data.AdditionalRecipients
	}

	return map[string]interface{}{
		"id":                         rule.Values["id"],
		"ruleType":                   rule.Values["ruleType"],
		"target":                     rule.Values["target"],
		"recipientType":              rule.Values["recipientType"],
		"notificationType":           rule.Values["notificationType"],
		"notificationLevel":          level,
		"isDefaultRecipientsEnabled": defaultRecipients,
		"notificationRecipients":     additionalRecipients,
	}
}

func flattenNotificationSettings(rule *map[string]interface{}) *RoleManagementPolicyNotificationSettings {
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
