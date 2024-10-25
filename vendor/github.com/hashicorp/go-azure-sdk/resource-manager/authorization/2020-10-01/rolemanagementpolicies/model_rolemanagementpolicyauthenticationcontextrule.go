package rolemanagementpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RoleManagementPolicyRule = RoleManagementPolicyAuthenticationContextRule{}

type RoleManagementPolicyAuthenticationContextRule struct {
	ClaimValue *string `json:"claimValue,omitempty"`
	IsEnabled  *bool   `json:"isEnabled,omitempty"`

	// Fields inherited from RoleManagementPolicyRule

	Id       *string                         `json:"id,omitempty"`
	RuleType RoleManagementPolicyRuleType    `json:"ruleType"`
	Target   *RoleManagementPolicyRuleTarget `json:"target,omitempty"`
}

func (s RoleManagementPolicyAuthenticationContextRule) RoleManagementPolicyRule() BaseRoleManagementPolicyRuleImpl {
	return BaseRoleManagementPolicyRuleImpl{
		Id:       s.Id,
		RuleType: s.RuleType,
		Target:   s.Target,
	}
}

var _ json.Marshaler = RoleManagementPolicyAuthenticationContextRule{}

func (s RoleManagementPolicyAuthenticationContextRule) MarshalJSON() ([]byte, error) {
	type wrapper RoleManagementPolicyAuthenticationContextRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RoleManagementPolicyAuthenticationContextRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RoleManagementPolicyAuthenticationContextRule: %+v", err)
	}

	decoded["ruleType"] = "RoleManagementPolicyAuthenticationContextRule"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RoleManagementPolicyAuthenticationContextRule: %+v", err)
	}

	return encoded, nil
}
