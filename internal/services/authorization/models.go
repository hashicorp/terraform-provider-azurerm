// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import "github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies"

var _ rolemanagementpolicies.RoleManagementPolicyRule = RoleManagementPolicyRule{}

// RoleManagementPolicyRule is a temporary workaround because the child models for rolemanagementpolicies.RoleManagementPolicyRule
// do not yet exist, due to issues resolving these child types in the Pandora importer.
// TODO: replace this type with the correct child models from the SDK package, once they are implemented
type RoleManagementPolicyRule struct {
	Id       *string                                             `json:"id,omitempty"`
	RuleType rolemanagementpolicies.RoleManagementPolicyRuleType `json:"ruleType"`

	ClaimValue                 *string                `json:"claimValue,omitempty"`
	EnabledRules               *[]string              `json:"enabledRules,omitempty"`
	IsDefaultRecipientsEnabled *bool                  `json:"isDefaultRecipientsEnabled"`
	IsEnabled                  *bool                  `json:"isEnabled,omitempty"`
	IsExpirationRequired       *bool                  `json:"isExpirationRequired,omitempty"`
	MaximumDuration            *string                `json:"maximumDuration,omitempty"`
	NotificationLevel          *string                `json:"notificationLevel"`
	NotificationRecipients     *[]string              `json:"notificationRecipients"`
	NotificationType           *string                `json:"notificationType"`
	RecipientType              *string                `json:"recipientType"`
	Setting                    map[string]interface{} `json:"setting,omitempty"`
	Target                     map[string]interface{} `json:"target,omitempty"`
}

func (r RoleManagementPolicyRule) RoleManagementPolicyRule() rolemanagementpolicies.BaseRoleManagementPolicyRuleImpl {
	return rolemanagementpolicies.BaseRoleManagementPolicyRuleImpl{
		Id:       r.Id,
		RuleType: r.RuleType,
	}
}
