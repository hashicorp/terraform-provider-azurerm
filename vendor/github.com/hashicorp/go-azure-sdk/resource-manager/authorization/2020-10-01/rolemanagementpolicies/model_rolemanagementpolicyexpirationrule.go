package rolemanagementpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RoleManagementPolicyRule = RoleManagementPolicyExpirationRule{}

type RoleManagementPolicyExpirationRule struct {
	IsExpirationRequired *bool   `json:"isExpirationRequired,omitempty"`
	MaximumDuration      *string `json:"maximumDuration,omitempty"`

	// Fields inherited from RoleManagementPolicyRule

	Id       *string                         `json:"id,omitempty"`
	RuleType RoleManagementPolicyRuleType    `json:"ruleType"`
	Target   *RoleManagementPolicyRuleTarget `json:"target,omitempty"`
}

func (s RoleManagementPolicyExpirationRule) RoleManagementPolicyRule() BaseRoleManagementPolicyRuleImpl {
	return BaseRoleManagementPolicyRuleImpl{
		Id:       s.Id,
		RuleType: s.RuleType,
		Target:   s.Target,
	}
}

var _ json.Marshaler = RoleManagementPolicyExpirationRule{}

func (s RoleManagementPolicyExpirationRule) MarshalJSON() ([]byte, error) {
	type wrapper RoleManagementPolicyExpirationRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RoleManagementPolicyExpirationRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RoleManagementPolicyExpirationRule: %+v", err)
	}

	decoded["ruleType"] = "RoleManagementPolicyExpirationRule"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RoleManagementPolicyExpirationRule: %+v", err)
	}

	return encoded, nil
}
