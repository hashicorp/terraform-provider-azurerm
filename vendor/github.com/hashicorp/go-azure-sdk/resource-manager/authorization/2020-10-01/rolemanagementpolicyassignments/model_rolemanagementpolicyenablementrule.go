package rolemanagementpolicyassignments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RoleManagementPolicyRule = RoleManagementPolicyEnablementRule{}

type RoleManagementPolicyEnablementRule struct {
	EnabledRules *[]EnablementRules `json:"enabledRules,omitempty"`

	// Fields inherited from RoleManagementPolicyRule

	Id       *string                         `json:"id,omitempty"`
	RuleType RoleManagementPolicyRuleType    `json:"ruleType"`
	Target   *RoleManagementPolicyRuleTarget `json:"target,omitempty"`
}

func (s RoleManagementPolicyEnablementRule) RoleManagementPolicyRule() BaseRoleManagementPolicyRuleImpl {
	return BaseRoleManagementPolicyRuleImpl{
		Id:       s.Id,
		RuleType: s.RuleType,
		Target:   s.Target,
	}
}

var _ json.Marshaler = RoleManagementPolicyEnablementRule{}

func (s RoleManagementPolicyEnablementRule) MarshalJSON() ([]byte, error) {
	type wrapper RoleManagementPolicyEnablementRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RoleManagementPolicyEnablementRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RoleManagementPolicyEnablementRule: %+v", err)
	}

	decoded["ruleType"] = "RoleManagementPolicyEnablementRule"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RoleManagementPolicyEnablementRule: %+v", err)
	}

	return encoded, nil
}
