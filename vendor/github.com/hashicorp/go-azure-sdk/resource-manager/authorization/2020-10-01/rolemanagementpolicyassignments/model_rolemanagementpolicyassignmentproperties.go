package rolemanagementpolicyassignments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleManagementPolicyAssignmentProperties struct {
	EffectiveRules             *[]RoleManagementPolicyRule `json:"effectiveRules,omitempty"`
	PolicyAssignmentProperties *PolicyAssignmentProperties `json:"policyAssignmentProperties,omitempty"`
	PolicyId                   *string                     `json:"policyId,omitempty"`
	RoleDefinitionId           *string                     `json:"roleDefinitionId,omitempty"`
	Scope                      *string                     `json:"scope,omitempty"`
}

var _ json.Unmarshaler = &RoleManagementPolicyAssignmentProperties{}

func (s *RoleManagementPolicyAssignmentProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		PolicyAssignmentProperties *PolicyAssignmentProperties `json:"policyAssignmentProperties,omitempty"`
		PolicyId                   *string                     `json:"policyId,omitempty"`
		RoleDefinitionId           *string                     `json:"roleDefinitionId,omitempty"`
		Scope                      *string                     `json:"scope,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.PolicyAssignmentProperties = decoded.PolicyAssignmentProperties
	s.PolicyId = decoded.PolicyId
	s.RoleDefinitionId = decoded.RoleDefinitionId
	s.Scope = decoded.Scope

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling RoleManagementPolicyAssignmentProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["effectiveRules"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling EffectiveRules into list []json.RawMessage: %+v", err)
		}

		output := make([]RoleManagementPolicyRule, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalRoleManagementPolicyRuleImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'EffectiveRules' for 'RoleManagementPolicyAssignmentProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.EffectiveRules = &output
	}

	return nil
}
