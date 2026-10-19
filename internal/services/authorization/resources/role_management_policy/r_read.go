// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_management_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Read() sdk.ResourceFunc {
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
						switch rule := r.(type) {
						case rolemanagementpolicies.RoleManagementPolicyAuthenticationContextRule:
							if rule.Id != nil && *rule.Id == "AuthenticationContext_EndUser_Assignment" {
								if rule.ClaimValue != nil && *rule.ClaimValue != "" {
									state.ActivationRules[0].RequireConditionalAccessContext = *rule.ClaimValue
								}
							}
						case rolemanagementpolicies.RoleManagementPolicyApprovalRule:
							if rule.Id != nil && *rule.Id == "Approval_EndUser_Assignment" {
								if rule.Setting != nil {
									settings := *rule.Setting
									state.ActivationRules[0].RequireApproval = pointer.From(settings.IsApprovalRequired)
									if settings.ApprovalStages != nil {
										approvalStages := *settings.ApprovalStages
										state.ActivationRules[0].ApprovalStages = make([]RoleManagementPolicyApprovalStage, 1)
										approvalStage := approvalStages[0]
										if primaryApprovers := approvalStage.PrimaryApprovers; primaryApprovers != nil && len(*primaryApprovers) > 0 {
											state.ActivationRules[0].ApprovalStages[0].PrimaryApprovers = make([]RoleManagementPolicyApprover, len(*primaryApprovers))
											for ia, pa := range *primaryApprovers {
												state.ActivationRules[0].ApprovalStages[0].PrimaryApprovers[ia] = RoleManagementPolicyApprover{
													ID:   pointer.From(pa.Id),
													Type: pointer.FromEnum(pa.UserType),
												}
											}
										}
									}
								}
							}
						case rolemanagementpolicies.RoleManagementPolicyEnablementRule:
							if rule.Id != nil && *rule.Id == "Enablement_Admin_Assignment" {
								state.ActiveAssignmentRules[0].RequireMultiFactorAuth = false
								state.ActiveAssignmentRules[0].RequireJustification = false

								if rule.EnabledRules != nil {
									for _, enabledRule := range *rule.EnabledRules {
										switch enabledRule {
										case rolemanagementpolicies.EnablementRulesMultiFactorAuthentication:
											state.ActiveAssignmentRules[0].RequireMultiFactorAuth = true
										case rolemanagementpolicies.EnablementRulesJustification:
											state.ActiveAssignmentRules[0].RequireJustification = true
										}
									}
								}
							}

							if rule.Id != nil && *rule.Id == "Enablement_EndUser_Assignment" {
								state.ActivationRules[0].RequireMultiFactorAuth = false
								state.ActivationRules[0].RequireJustification = false
								state.ActivationRules[0].RequireTicketInfo = false

								if rule.EnabledRules != nil {
									for _, enabledRule := range *rule.EnabledRules {
										switch enabledRule {
										case rolemanagementpolicies.EnablementRulesMultiFactorAuthentication:
											state.ActivationRules[0].RequireMultiFactorAuth = true
										case rolemanagementpolicies.EnablementRulesJustification:
											state.ActivationRules[0].RequireJustification = true
										case rolemanagementpolicies.EnablementRulesTicketing:
											state.ActivationRules[0].RequireTicketInfo = true
										}
									}
								}
							}

						case rolemanagementpolicies.RoleManagementPolicyExpirationRule:
							if rule.Id != nil {
								switch *rule.Id {
								case "Expiration_Admin_Eligibility":
									state.EligibleAssignmentRules[0].ExpirationRequired = pointer.From(rule.IsExpirationRequired)
									state.EligibleAssignmentRules[0].ExpireAfter = pointer.From(rule.MaximumDuration)

								case "Expiration_Admin_Assignment":
									state.ActiveAssignmentRules[0].ExpirationRequired = pointer.From(rule.IsExpirationRequired)
									state.ActiveAssignmentRules[0].ExpireAfter = pointer.From(rule.MaximumDuration)

								case "Expiration_EndUser_Assignment":
									state.ActivationRules[0].MaximumDuration = pointer.From(rule.MaximumDuration)
								}
							}
						case rolemanagementpolicies.RoleManagementPolicyNotificationRule:
							if rule.Id != nil {
								switch *rule.Id {
								case "Notification_Admin_Admin_Assignment":
									state.NotificationRules[0].ActiveAssignments[0].AdminNotifications = []RoleManagementPolicyNotificationSettings{
										*flattenNotificationSettings(rule),
									}

								case "Notification_Admin_Admin_Eligibility":
									state.NotificationRules[0].EligibleAssignments[0].AdminNotifications = []RoleManagementPolicyNotificationSettings{
										*flattenNotificationSettings(rule),
									}

								case "Notification_Admin_EndUser_Assignment":
									state.NotificationRules[0].EligibleActivations[0].AdminNotifications = []RoleManagementPolicyNotificationSettings{
										*flattenNotificationSettings(rule),
									}

								case "Notification_Approver_Admin_Assignment":
									state.NotificationRules[0].ActiveAssignments[0].ApproverNotifications = []RoleManagementPolicyNotificationSettings{
										*flattenNotificationSettings(rule),
									}

								case "Notification_Approver_Admin_Eligibility":
									state.NotificationRules[0].EligibleAssignments[0].ApproverNotifications = []RoleManagementPolicyNotificationSettings{
										*flattenNotificationSettings(rule),
									}

								case "Notification_Approver_EndUser_Assignment":
									state.NotificationRules[0].EligibleActivations[0].ApproverNotifications = []RoleManagementPolicyNotificationSettings{
										*flattenNotificationSettings(rule),
									}

								case "Notification_Requestor_Admin_Assignment":
									state.NotificationRules[0].ActiveAssignments[0].AssigneeNotifications = []RoleManagementPolicyNotificationSettings{
										*flattenNotificationSettings(rule),
									}

								case "Notification_Requestor_Admin_Eligibility":
									state.NotificationRules[0].EligibleAssignments[0].AssigneeNotifications = []RoleManagementPolicyNotificationSettings{
										*flattenNotificationSettings(rule),
									}

								case "Notification_Requestor_EndUser_Assignment":
									state.NotificationRules[0].EligibleActivations[0].AssigneeNotifications = []RoleManagementPolicyNotificationSettings{
										*flattenNotificationSettings(rule),
									}
								}
							}
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
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

func notificationRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"admin_notifications": {
			Description: "Admin notification settings",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: notificationSettingsSchema(),
			},
		},

		"approver_notifications": {
			Description: "Approver notification settings",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: notificationSettingsSchema(),
			},
		},

		"assignee_notifications": {
			Description: "Assignee notification settings",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
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
			Computed:    true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
