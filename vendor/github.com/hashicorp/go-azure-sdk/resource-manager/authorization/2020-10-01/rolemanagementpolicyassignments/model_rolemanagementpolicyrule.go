package rolemanagementpolicyassignments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleManagementPolicyRule interface {
	RoleManagementPolicyRule() BaseRoleManagementPolicyRuleImpl
}

var _ RoleManagementPolicyRule = BaseRoleManagementPolicyRuleImpl{}

type BaseRoleManagementPolicyRuleImpl struct {
	Id       *string                         `json:"id,omitempty"`
	RuleType RoleManagementPolicyRuleType    `json:"ruleType"`
	Target   *RoleManagementPolicyRuleTarget `json:"target,omitempty"`
}

func (s BaseRoleManagementPolicyRuleImpl) RoleManagementPolicyRule() BaseRoleManagementPolicyRuleImpl {
	return s
}

var _ RoleManagementPolicyRule = RawRoleManagementPolicyRuleImpl{}

// RawRoleManagementPolicyRuleImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawRoleManagementPolicyRuleImpl struct {
	roleManagementPolicyRule BaseRoleManagementPolicyRuleImpl
	Type                     string
	Values                   map[string]interface{}
}

func (s RawRoleManagementPolicyRuleImpl) RoleManagementPolicyRule() BaseRoleManagementPolicyRuleImpl {
	return s.roleManagementPolicyRule
}

func UnmarshalRoleManagementPolicyRuleImplementation(input []byte) (RoleManagementPolicyRule, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RoleManagementPolicyRule into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["ruleType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	var parent BaseRoleManagementPolicyRuleImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseRoleManagementPolicyRuleImpl: %+v", err)
	}

	return RawRoleManagementPolicyRuleImpl{
		roleManagementPolicyRule: parent,
		Type:                     value,
		Values:                   temp,
	}, nil

}
