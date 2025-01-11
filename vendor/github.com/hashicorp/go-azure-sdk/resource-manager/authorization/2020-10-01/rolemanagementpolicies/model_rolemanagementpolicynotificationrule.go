package rolemanagementpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RoleManagementPolicyRule = RoleManagementPolicyNotificationRule{}

type RoleManagementPolicyNotificationRule struct {
	IsDefaultRecipientsEnabled *bool                          `json:"isDefaultRecipientsEnabled,omitempty"`
	NotificationLevel          *NotificationLevel             `json:"notificationLevel,omitempty"`
	NotificationRecipients     *[]string                      `json:"notificationRecipients,omitempty"`
	NotificationType           *NotificationDeliveryMechanism `json:"notificationType,omitempty"`
	RecipientType              *RecipientType                 `json:"recipientType,omitempty"`

	// Fields inherited from RoleManagementPolicyRule

	Id       *string                         `json:"id,omitempty"`
	RuleType RoleManagementPolicyRuleType    `json:"ruleType"`
	Target   *RoleManagementPolicyRuleTarget `json:"target,omitempty"`
}

func (s RoleManagementPolicyNotificationRule) RoleManagementPolicyRule() BaseRoleManagementPolicyRuleImpl {
	return BaseRoleManagementPolicyRuleImpl{
		Id:       s.Id,
		RuleType: s.RuleType,
		Target:   s.Target,
	}
}

var _ json.Marshaler = RoleManagementPolicyNotificationRule{}

func (s RoleManagementPolicyNotificationRule) MarshalJSON() ([]byte, error) {
	type wrapper RoleManagementPolicyNotificationRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RoleManagementPolicyNotificationRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RoleManagementPolicyNotificationRule: %+v", err)
	}

	decoded["ruleType"] = "RoleManagementPolicyNotificationRule"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RoleManagementPolicyNotificationRule: %+v", err)
	}

	return encoded, nil
}
