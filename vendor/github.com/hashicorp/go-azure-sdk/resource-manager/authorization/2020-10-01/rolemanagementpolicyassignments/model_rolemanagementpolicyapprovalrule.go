package rolemanagementpolicyassignments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RoleManagementPolicyRule = RoleManagementPolicyApprovalRule{}

type RoleManagementPolicyApprovalRule struct {
	Setting *ApprovalSettings `json:"setting,omitempty"`

	// Fields inherited from RoleManagementPolicyRule

	Id       *string                         `json:"id,omitempty"`
	RuleType RoleManagementPolicyRuleType    `json:"ruleType"`
	Target   *RoleManagementPolicyRuleTarget `json:"target,omitempty"`
}

func (s RoleManagementPolicyApprovalRule) RoleManagementPolicyRule() BaseRoleManagementPolicyRuleImpl {
	return BaseRoleManagementPolicyRuleImpl{
		Id:       s.Id,
		RuleType: s.RuleType,
		Target:   s.Target,
	}
}

var _ json.Marshaler = RoleManagementPolicyApprovalRule{}

func (s RoleManagementPolicyApprovalRule) MarshalJSON() ([]byte, error) {
	type wrapper RoleManagementPolicyApprovalRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RoleManagementPolicyApprovalRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RoleManagementPolicyApprovalRule: %+v", err)
	}

	decoded["ruleType"] = "RoleManagementPolicyApprovalRule"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RoleManagementPolicyApprovalRule: %+v", err)
	}

	return encoded, nil
}
