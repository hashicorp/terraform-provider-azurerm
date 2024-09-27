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
	if len(*resp.Model) != 1 {
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

	// Take the slice of rules and convert it to a map with the ID as the key
	existingRules := make(map[string]rolemanagementpolicies.RawRoleManagementPolicyRuleImpl)
	if props := rolePolicy.Properties; props != nil {
		if props.Rules != nil {
			for _, r := range *rolePolicy.Properties.Rules { // TODO
				rule := r.(rolemanagementpolicies.RawRoleManagementPolicyRuleImpl)
				existingRules[rule.Values["id"].(string)] = rule
			}
		}
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

			updatedRules = append(updatedRules, RoleManagementPolicyRule{
				Id:                   pointer.To(id),
				RuleType:             rolemanagementpolicies.RoleManagementPolicyRuleType(ruleType),
				Target:               target,
				IsExpirationRequired: pointer.To(expirationRequired),
				MaximumDuration:      pointer.To(maximumDuration),
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

			updatedRules = append(updatedRules, RoleManagementPolicyRule{
				Id:           pointer.To(id),
				RuleType:     rolemanagementpolicies.RoleManagementPolicyRuleType(ruleType),
				Target:       target,
				EnabledRules: pointer.To(enabledRules),
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

			updatedRules = append(updatedRules, RoleManagementPolicyRule{
				Id:                   pointer.To(id),
				RuleType:             rolemanagementpolicies.RoleManagementPolicyRuleType(ruleType),
				Target:               target,
				IsExpirationRequired: pointer.To(expirationRequired),
				MaximumDuration:      pointer.To(maximumDuration),
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

			updatedRules = append(updatedRules, RoleManagementPolicyRule{
				Id:              pointer.To(id),
				RuleType:        rolemanagementpolicies.RoleManagementPolicyRuleType(ruleType),
				Target:          target,
				MaximumDuration: pointer.To(maximumDuration),
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
						for _, stage := range approvalStagesRaw.([]interface{}) {
							if v, ok := stage.(map[string]interface{}); ok {
								approvalStages = append(approvalStages, v)
							}
						}
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

			updatedRules = append(updatedRules, RoleManagementPolicyRule{
				Id:       pointer.To(id),
				RuleType: rolemanagementpolicies.RoleManagementPolicyRuleType(ruleType),
				Target:   target,
				Setting: map[string]interface{}{
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
			if _, ok = metadata.ResourceData.GetOk("activation_rules.0.required_conditional_access_authentication_context"); ok {
				isEnabled = true
				if len(model.ActivationRules) == 1 {
					claimValue = model.ActivationRules[0].RequireConditionalAccessContext
				}
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

			updatedRules = append(updatedRules, RoleManagementPolicyRule{
				Id:         pointer.To(id),
				RuleType:   rolemanagementpolicies.RoleManagementPolicyRuleType(ruleType),
				Target:     target,
				IsEnabled:  pointer.To(isEnabled),
				ClaimValue: pointer.To(claimValue),
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

			updatedRules = append(updatedRules, RoleManagementPolicyRule{
				Id:           pointer.To(id),
				RuleType:     rolemanagementpolicies.RoleManagementPolicyRuleType(ruleType),
				Target:       target,
				EnabledRules: pointer.To(enabledRules),
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

	return RoleManagementPolicyRule{
		Id:                         pointer.To(id),
		RuleType:                   rolemanagementpolicies.RoleManagementPolicyRuleType(ruleType),
		Target:                     target,
		RecipientType:              pointer.To(recipientType),
		NotificationType:           pointer.To(notificationType),
		NotificationLevel:          pointer.To(level),
		IsDefaultRecipientsEnabled: pointer.To(defaultRecipients),
		NotificationRecipients:     pointer.To(additionalRecipients),
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
