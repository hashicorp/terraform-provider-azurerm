// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_management_policy

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

func (r Resource) ModelObject() any {
	return &RoleManagementPolicyModel{}
}
