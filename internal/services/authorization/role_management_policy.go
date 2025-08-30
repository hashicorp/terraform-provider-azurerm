// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func FindRoleManagementPolicyId(ctx context.Context, client *rolemanagementpolicies.RoleManagementPoliciesClient, scope string, roleDefinitionId string) (*rolemanagementpolicies.ScopedRoleManagementPolicyId, error) {
	// List role management policies to find the latest policy for the provided scope and role definition
	options := rolemanagementpolicies.ListForScopeOperationOptions{
		Filter: pointer.To(fmt.Sprintf("roleDefinitionId eq '%s'", odata.EscapeSingleQuote(roleDefinitionId))),
	}
	scopeId, err := commonids.ParseScopeID(scope)
	if err != nil {
		return nil, err
	}
	resp, err := client.ListForScope(ctx, *scopeId, options)
	if err != nil {
		return nil, fmt.Errorf("listing Role Management Policies for %s and Role Definition %q: %+v", scope, roleDefinitionId, err)
	}

	// There should be one policy to represent a given scope and role definition
	if resp.Model == nil {
		return nil, fmt.Errorf("listing Role Management Policies for %s and Role Definition %q: result was nil", scope, roleDefinitionId)
	}

	if len(*resp.Model) == 0 {
		return nil, fmt.Errorf("no Role Management Policy returned for %s and Role Definition %q", scope, roleDefinitionId)
	}

	if len(*resp.Model) > 1 {
		return nil, fmt.Errorf("more than one Role Management Policy returned for %s and Role Definition %q", scope, roleDefinitionId)
	}

	policy := (*resp.Model)[0]
	if policy.Name == nil {
		return nil, fmt.Errorf("retrieving Role Management Policy for %s and Role Definition %q: `name` was nil", scope, roleDefinitionId)
	}

	// Note: the "Name" is actually a UUID that changes each time the policy is updated
	id := rolemanagementpolicies.NewScopedRoleManagementPolicyID(scope, *policy.Name)

	return &id, nil
}

func buildRoleManagementPolicyForUpdate(metadata *sdk.ResourceMetaData, rolePolicy *rolemanagementpolicies.RoleManagementPolicy) (*rolemanagementpolicies.RoleManagementPolicy, error) {
	if rolePolicy == nil {
		return nil, fmt.Errorf("existing Role Management Policy was nil")
	}

	var model RoleManagementPolicyModel
	if err := metadata.Decode(&model); err != nil {
		return nil, fmt.Errorf("decoding: %+v", err)
	}

	return buildRoleManagementPolicyInternal(metadata, rolePolicy, &model)
}

func buildRoleManagementPolicyForDelete(rolePolicy *rolemanagementpolicies.RoleManagementPolicy) (*rolemanagementpolicies.RoleManagementPolicy, error) {
	if rolePolicy == nil {
		return nil, fmt.Errorf("existing Role Management Policy was nil")
	}

	defaultsModel := buildModelWithDefaultNotificationSettings()
	return buildRoleManagementPolicyInternal(nil, rolePolicy, defaultsModel)
}

func buildModelWithDefaultNotificationSettings() *RoleManagementPolicyModel {
	return &RoleManagementPolicyModel{
		NotificationRules: []RoleManagementPolicyNotificationEvents{
			{
				ActiveAssignments: []RoleManagementPolicyNotificationRule{
					{
						AdminNotifications: []RoleManagementPolicyNotificationSettings{
							{
								NotificationLevel:    "All",
								DefaultRecipients:    true,
								AdditionalRecipients: []string{},
							},
						},
						ApproverNotifications: []RoleManagementPolicyNotificationSettings{
							{
								NotificationLevel:    "All",
								DefaultRecipients:    true,
								AdditionalRecipients: []string{},
							},
						},
						AssigneeNotifications: []RoleManagementPolicyNotificationSettings{
							{
								NotificationLevel:    "All",
								DefaultRecipients:    true,
								AdditionalRecipients: []string{},
							},
						},
					},
				},
				EligibleActivations: []RoleManagementPolicyNotificationRule{
					{
						AdminNotifications: []RoleManagementPolicyNotificationSettings{
							{
								NotificationLevel:    "All",
								DefaultRecipients:    true,
								AdditionalRecipients: []string{},
							},
						},
						ApproverNotifications: []RoleManagementPolicyNotificationSettings{
							{
								NotificationLevel:    "All",
								DefaultRecipients:    true,
								AdditionalRecipients: []string{},
							},
						},
						AssigneeNotifications: []RoleManagementPolicyNotificationSettings{
							{
								NotificationLevel:    "All",
								DefaultRecipients:    true,
								AdditionalRecipients: []string{},
							},
						},
					},
				},
				EligibleAssignments: []RoleManagementPolicyNotificationRule{
					{
						AdminNotifications: []RoleManagementPolicyNotificationSettings{
							{
								NotificationLevel:    "All",
								DefaultRecipients:    true,
								AdditionalRecipients: []string{},
							},
						},
						ApproverNotifications: []RoleManagementPolicyNotificationSettings{
							{
								NotificationLevel:    "All",
								DefaultRecipients:    true,
								AdditionalRecipients: []string{},
							},
						},
						AssigneeNotifications: []RoleManagementPolicyNotificationSettings{
							{
								NotificationLevel:    "All",
								DefaultRecipients:    true,
								AdditionalRecipients: []string{},
							},
						},
					},
				},
			},
		},
	}
}

func buildRoleManagementPolicyInternal(metadata *sdk.ResourceMetaData, rolePolicy *rolemanagementpolicies.RoleManagementPolicy, model *RoleManagementPolicyModel) (*rolemanagementpolicies.RoleManagementPolicy, error) {
	existingRules := make(map[string]rolemanagementpolicies.RoleManagementPolicyRule)
	if props := rolePolicy.Properties; props != nil {
		if props.Rules != nil {
			for _, r := range *rolePolicy.Properties.Rules {
				if id := pointer.From(r.RoleManagementPolicyRule().Id); id != "" {
					existingRules[id] = r
				}
			}
		}
	}
	updatedRules := make([]rolemanagementpolicies.RoleManagementPolicyRule, 0)

	if metadata == nil || metadata.ResourceData.HasChange("eligible_assignment_rules") {
		if expirationAdminEligibilityBase, ok := existingRules["Expiration_Admin_Eligibility"]; ok {
			if expirationAdminEligibility, ok := expirationAdminEligibilityBase.(rolemanagementpolicies.RoleManagementPolicyExpirationRule); ok {
				expirationRequired := pointer.From(expirationAdminEligibility.IsExpirationRequired)
				maximumDuration := pointer.From(expirationAdminEligibility.MaximumDuration)

				if len(model.EligibleAssignmentRules) == 1 {
					if expirationRequired != model.EligibleAssignmentRules[0].ExpirationRequired {
						expirationAdminEligibility.IsExpirationRequired = pointer.To(model.EligibleAssignmentRules[0].ExpirationRequired)
					}
					if maximumDuration != model.EligibleAssignmentRules[0].ExpireAfter &&
						model.EligibleAssignmentRules[0].ExpireAfter != "" {
						expirationAdminEligibility.MaximumDuration = pointer.To(model.EligibleAssignmentRules[0].ExpireAfter)
					}
				}

				updatedRules = append(updatedRules, expirationAdminEligibility)
			}
		}
	}

	if metadata == nil || metadata.ResourceData.HasChange("active_assignment_rules.0.require_multifactor_authentication") ||
		metadata.ResourceData.HasChange("active_assignment_rules.0.require_justification") ||
		metadata.ResourceData.HasChange("active_assignment_rules.0.require_ticket_info") {
		if enablementAdminEligibilityBase, ok := existingRules["Enablement_Admin_Assignment"]; ok {
			if enablementAdminEligibility, ok := enablementAdminEligibilityBase.(rolemanagementpolicies.RoleManagementPolicyEnablementRule); ok {
				enabledRules := make([]rolemanagementpolicies.EnablementRules, 0)
				if len(model.ActiveAssignmentRules) == 1 {
					if model.ActiveAssignmentRules[0].RequireMultiFactorAuth {
						enabledRules = append(enabledRules, rolemanagementpolicies.EnablementRulesMultiFactorAuthentication)
					}
					if model.ActiveAssignmentRules[0].RequireJustification {
						enabledRules = append(enabledRules, rolemanagementpolicies.EnablementRulesJustification)
					}
					if model.ActiveAssignmentRules[0].RequireTicketInfo {
						enabledRules = append(enabledRules, rolemanagementpolicies.EnablementRulesTicketing)
					}
				}
				enablementAdminEligibility.EnabledRules = pointer.To(enabledRules)
				updatedRules = append(updatedRules, enablementAdminEligibility)
			}
		}
	}

	if metadata == nil || metadata.ResourceData.HasChange("active_assignment_rules.0.expiration_required") ||
		metadata.ResourceData.HasChange("active_assignment_rules.0.expire_after") {
		if expirationAdminAssignmentBase, ok := existingRules["Expiration_Admin_Assignment"]; ok {
			if expirationAdminAssignment, ok := expirationAdminAssignmentBase.(rolemanagementpolicies.RoleManagementPolicyExpirationRule); ok {
				expirationRequired := pointer.From(expirationAdminAssignment.IsExpirationRequired)
				maximumDuration := pointer.From(expirationAdminAssignment.MaximumDuration)

				if len(model.ActiveAssignmentRules) == 1 {
					if expirationRequired != model.ActiveAssignmentRules[0].ExpirationRequired {
						expirationAdminAssignment.IsExpirationRequired = pointer.To(model.ActiveAssignmentRules[0].ExpirationRequired)
					}
					if maximumDuration != model.ActiveAssignmentRules[0].ExpireAfter &&
						model.ActiveAssignmentRules[0].ExpireAfter != "" {
						expirationAdminAssignment.MaximumDuration = pointer.To(model.ActiveAssignmentRules[0].ExpireAfter)
					}
				}

				updatedRules = append(updatedRules, expirationAdminAssignment)
			}
		}
	}

	if metadata == nil || metadata.ResourceData.HasChange("activation_rules.0.maximum_duration") {
		if expirationEndUserAssignmentBase, ok := existingRules["Expiration_EndUser_Assignment"]; ok {
			if expirationEndUserAssignment, ok := expirationEndUserAssignmentBase.(rolemanagementpolicies.RoleManagementPolicyExpirationRule); ok {
				if len(model.ActivationRules) == 1 {
					expirationEndUserAssignment.MaximumDuration = pointer.To(model.ActivationRules[0].MaximumDuration)
				}

				updatedRules = append(updatedRules, expirationEndUserAssignment)
			}
		}
	}

	if metadata == nil || metadata.ResourceData.HasChange("activation_rules.0.require_approval") ||
		metadata.ResourceData.HasChange("activation_rules.0.approval_stage") {
		if approvalEndUserAssignmentBase, ok := existingRules["Approval_EndUser_Assignment"]; ok {
			if approvalEndUserAssignment, ok := approvalEndUserAssignmentBase.(rolemanagementpolicies.RoleManagementPolicyApprovalRule); ok {
				if len(model.ActivationRules) == 1 {
					if model.ActivationRules[0].RequireApproval && len(model.ActivationRules[0].ApprovalStages) != 1 {
						return nil, fmt.Errorf("require_approval is true, but no approval_stages are provided")
					}
				}

				if settings := approvalEndUserAssignment.Setting; settings != nil {
					if len(model.ActivationRules) == 1 {
						if pointer.From(settings.IsApprovalRequired) != model.ActivationRules[0].RequireApproval {
							settings.IsApprovalRequired = pointer.To(model.ActivationRules[0].RequireApproval)
						}
					}

					if metadata == nil || metadata.ResourceData.HasChange("activation_rules.0.approval_stage") {
						if len(model.ActivationRules) == 1 {
							approvalStages := make([]rolemanagementpolicies.ApprovalStage, len(model.ActivationRules[0].ApprovalStages))
							for i, stage := range model.ActivationRules[0].ApprovalStages {
								primaryApprovers := make([]rolemanagementpolicies.UserSet, len(stage.PrimaryApprovers))
								for ia, approver := range stage.PrimaryApprovers {
									primaryApprovers[ia] = rolemanagementpolicies.UserSet{
										Id:       pointer.To(approver.ID),
										UserType: pointer.To(rolemanagementpolicies.UserType(approver.Type)),
									}
								}

								approvalStages[i] = rolemanagementpolicies.ApprovalStage{
									PrimaryApprovers: &primaryApprovers,
								}
							}
							settings.ApprovalStages = &approvalStages
						}
					}
				}

				updatedRules = append(updatedRules, approvalEndUserAssignment)
			}
		}
	}

	if metadata == nil || metadata.ResourceData.HasChange("activation_rules.0.required_conditional_access_authentication_context") {
		if authEndUserAssignmentBase, ok := existingRules["AuthenticationContext_EndUser_Assignment"]; ok {
			if authEndUserAssignment, ok := authEndUserAssignmentBase.(rolemanagementpolicies.RoleManagementPolicyAuthenticationContextRule); ok {
				if metadata == nil {
					authEndUserAssignment.IsEnabled = pointer.To(false)
					authEndUserAssignment.ClaimValue = pointer.To("")
				} else if _, ok := metadata.ResourceData.GetOk("activation_rules.0.required_conditional_access_authentication_context"); ok {
					authEndUserAssignment.IsEnabled = pointer.To(true)
					if len(model.ActivationRules) == 1 {
						authEndUserAssignment.ClaimValue = pointer.To(model.ActivationRules[0].RequireConditionalAccessContext)
					}
				}

				updatedRules = append(updatedRules, authEndUserAssignment)
			}
		}
	}

	if metadata == nil || metadata.ResourceData.HasChange("activation_rules.0.require_multifactor_authentication") ||
		metadata.ResourceData.HasChange("activation_rules.0.require_justification") ||
		metadata.ResourceData.HasChange("activation_rules.0.require_ticket_info") {
		if enablementEndUserAssignmentBase, ok := existingRules["Enablement_EndUser_Assignment"]; ok {
			if enablementEndUserAssignment, ok := enablementEndUserAssignmentBase.(rolemanagementpolicies.RoleManagementPolicyEnablementRule); ok {
				enabledRules := make([]rolemanagementpolicies.EnablementRules, 0)
				if len(model.ActivationRules) == 1 {
					if model.ActivationRules[0].RequireMultiFactorAuth {
						enabledRules = append(enabledRules, rolemanagementpolicies.EnablementRulesMultiFactorAuthentication)
					}
					if model.ActivationRules[0].RequireJustification {
						enabledRules = append(enabledRules, rolemanagementpolicies.EnablementRulesJustification)
					}
					if model.ActivationRules[0].RequireTicketInfo {
						enabledRules = append(enabledRules, rolemanagementpolicies.EnablementRulesTicketing)
					}
				}
				enablementEndUserAssignment.EnabledRules = &enabledRules

				updatedRules = append(updatedRules, enablementEndUserAssignment)
			}
		}
	}

	notificationRulesConfig := []struct {
		changeDetectionKey string
		recipientChangeKey string
		azureRuleID        string
		getModelSettings   func() *RoleManagementPolicyNotificationSettings
	}{
		{"notification_rules.0.eligible_assignments.0.admin_notifications", "notification_rules.0.eligible_assignments.0.admin_notifications.0.additional_recipients", "Notification_Admin_Admin_Eligibility", func() *RoleManagementPolicyNotificationSettings {
			if len(model.NotificationRules) == 1 && len(model.NotificationRules[0].EligibleAssignments) == 1 && len(model.NotificationRules[0].EligibleAssignments[0].AdminNotifications) == 1 {
				return &model.NotificationRules[0].EligibleAssignments[0].AdminNotifications[0]
			}
			return nil
		}},
		{"notification_rules.0.active_assignments.0.admin_notifications", "notification_rules.0.active_assignments.0.admin_notifications.0.additional_recipients", "Notification_Admin_Admin_Assignment", func() *RoleManagementPolicyNotificationSettings {
			if len(model.NotificationRules) == 1 && len(model.NotificationRules[0].ActiveAssignments) == 1 && len(model.NotificationRules[0].ActiveAssignments[0].AdminNotifications) == 1 {
				return &model.NotificationRules[0].ActiveAssignments[0].AdminNotifications[0]
			}
			return nil
		}},
		{"notification_rules.0.eligible_activations.0.admin_notifications", "notification_rules.0.eligible_activations.0.admin_notifications.0.additional_recipients", "Notification_Admin_EndUser_Assignment", func() *RoleManagementPolicyNotificationSettings {
			if len(model.NotificationRules) == 1 && len(model.NotificationRules[0].EligibleActivations) == 1 && len(model.NotificationRules[0].EligibleActivations[0].AdminNotifications) == 1 {
				return &model.NotificationRules[0].EligibleActivations[0].AdminNotifications[0]
			}
			return nil
		}},
		{"notification_rules.0.eligible_assignments.0.approver_notifications", "notification_rules.0.eligible_assignments.0.approver_notifications.0.additional_recipients", "Notification_Approver_Admin_Eligibility", func() *RoleManagementPolicyNotificationSettings {
			if len(model.NotificationRules) == 1 && len(model.NotificationRules[0].EligibleAssignments) == 1 && len(model.NotificationRules[0].EligibleAssignments[0].ApproverNotifications) == 1 {
				return &model.NotificationRules[0].EligibleAssignments[0].ApproverNotifications[0]
			}
			return nil
		}},
		{"notification_rules.0.active_assignments.0.approver_notifications", "notification_rules.0.active_assignments.0.approver_notifications.0.additional_recipients", "Notification_Approver_Admin_Assignment", func() *RoleManagementPolicyNotificationSettings {
			if len(model.NotificationRules) == 1 && len(model.NotificationRules[0].ActiveAssignments) == 1 && len(model.NotificationRules[0].ActiveAssignments[0].ApproverNotifications) == 1 {
				return &model.NotificationRules[0].ActiveAssignments[0].ApproverNotifications[0]
			}
			return nil
		}},
		{"notification_rules.0.eligible_activations.0.approver_notifications", "notification_rules.0.eligible_activations.0.approver_notifications.0.additional_recipients", "Notification_Approver_EndUser_Assignment", func() *RoleManagementPolicyNotificationSettings {
			if len(model.NotificationRules) == 1 && len(model.NotificationRules[0].EligibleActivations) == 1 && len(model.NotificationRules[0].EligibleActivations[0].ApproverNotifications) == 1 {
				return &model.NotificationRules[0].EligibleActivations[0].ApproverNotifications[0]
			}
			return nil
		}},
		{"notification_rules.0.eligible_assignments.0.assignee_notifications", "notification_rules.0.eligible_assignments.0.assignee_notifications.0.additional_recipients", "Notification_Requestor_Admin_Eligibility", func() *RoleManagementPolicyNotificationSettings {
			if len(model.NotificationRules) == 1 && len(model.NotificationRules[0].EligibleAssignments) == 1 && len(model.NotificationRules[0].EligibleAssignments[0].AssigneeNotifications) == 1 {
				return &model.NotificationRules[0].EligibleAssignments[0].AssigneeNotifications[0]
			}
			return nil
		}},
		{"notification_rules.0.active_assignments.0.assignee_notifications", "notification_rules.0.active_assignments.0.assignee_notifications.0.additional_recipients", "Notification_Requestor_Admin_Assignment", func() *RoleManagementPolicyNotificationSettings {
			if len(model.NotificationRules) == 1 && len(model.NotificationRules[0].ActiveAssignments) == 1 && len(model.NotificationRules[0].ActiveAssignments[0].AssigneeNotifications) == 1 {
				return &model.NotificationRules[0].ActiveAssignments[0].AssigneeNotifications[0]
			}
			return nil
		}},
		{"notification_rules.0.eligible_activations.0.assignee_notifications", "notification_rules.0.eligible_activations.0.assignee_notifications.0.additional_recipients", "Notification_Requestor_EndUser_Assignment", func() *RoleManagementPolicyNotificationSettings {
			if len(model.NotificationRules) == 1 && len(model.NotificationRules[0].EligibleActivations) == 1 && len(model.NotificationRules[0].EligibleActivations[0].AssigneeNotifications) == 1 {
				return &model.NotificationRules[0].EligibleActivations[0].AssigneeNotifications[0]
			}
			return nil
		}},
	}

	for _, ruleConfig := range notificationRulesConfig {
		if metadata == nil || metadata.ResourceData.HasChange(ruleConfig.changeDetectionKey) {
			if existingRuleBase, ok := existingRules[ruleConfig.azureRuleID]; ok {
				if existingRule, ok := existingRuleBase.(rolemanagementpolicies.RoleManagementPolicyNotificationRule); ok {
					configSettings := ruleConfig.getModelSettings()
					if configSettings != nil {
						recipientChange := metadata == nil || metadata.ResourceData.HasChange(ruleConfig.recipientChangeKey)
						updatedRules = append(updatedRules,
							expandNotificationSettings(
								existingRule,
								*configSettings,
								recipientChange,
							),
						)
					} else {
						defaultSettings := RoleManagementPolicyNotificationSettings{
							NotificationLevel:    "All",
							DefaultRecipients:    true,
							AdditionalRecipients: []string{},
						}
						updatedRules = append(updatedRules,
							expandNotificationSettings(existingRule, defaultSettings, true),
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

func expandNotificationSettings(rule rolemanagementpolicies.RoleManagementPolicyNotificationRule, data RoleManagementPolicyNotificationSettings, recipientChange bool) rolemanagementpolicies.RoleManagementPolicyRule {
	if pointer.From(rule.NotificationLevel) != rolemanagementpolicies.NotificationLevel(data.NotificationLevel) {
		rule.NotificationLevel = pointer.To(rolemanagementpolicies.NotificationLevel(data.NotificationLevel))
	}

	if pointer.From(rule.IsDefaultRecipientsEnabled) != data.DefaultRecipients {
		rule.IsDefaultRecipientsEnabled = pointer.To(data.DefaultRecipients)
	}

	if recipientChange {
		rule.NotificationRecipients = pointer.To(data.AdditionalRecipients)
	}

	return rule
}

func flattenNotificationSettings(rule rolemanagementpolicies.RoleManagementPolicyNotificationRule) *RoleManagementPolicyNotificationSettings {
	return &RoleManagementPolicyNotificationSettings{
		NotificationLevel:    string(pointer.From(rule.NotificationLevel)),
		DefaultRecipients:    pointer.From(rule.IsDefaultRecipientsEnabled),
		AdditionalRecipients: pointer.From(rule.NotificationRecipients),
	}
}
