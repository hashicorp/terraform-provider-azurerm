package firewallpolicyrulecollectiongroups

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FirewallPolicyRuleCollection = FirewallPolicyNatRuleCollection{}

type FirewallPolicyNatRuleCollection struct {
	Action *FirewallPolicyNatRuleCollectionAction `json:"action,omitempty"`
	Rules  *[]FirewallPolicyRule                  `json:"rules,omitempty"`

	// Fields inherited from FirewallPolicyRuleCollection
	Name     *string `json:"name,omitempty"`
	Priority *int64  `json:"priority,omitempty"`
}

var _ json.Marshaler = FirewallPolicyNatRuleCollection{}

func (s FirewallPolicyNatRuleCollection) MarshalJSON() ([]byte, error) {
	type wrapper FirewallPolicyNatRuleCollection
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FirewallPolicyNatRuleCollection: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FirewallPolicyNatRuleCollection: %+v", err)
	}
	decoded["ruleCollectionType"] = "FirewallPolicyNatRuleCollection"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FirewallPolicyNatRuleCollection: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &FirewallPolicyNatRuleCollection{}

func (s *FirewallPolicyNatRuleCollection) UnmarshalJSON(bytes []byte) error {
	type alias FirewallPolicyNatRuleCollection
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into FirewallPolicyNatRuleCollection: %+v", err)
	}

	s.Action = decoded.Action
	s.Name = decoded.Name
	s.Priority = decoded.Priority

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling FirewallPolicyNatRuleCollection into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["rules"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Rules into list []json.RawMessage: %+v", err)
		}

		output := make([]FirewallPolicyRule, 0)
		for i, val := range listTemp {
			impl, err := unmarshalFirewallPolicyRuleImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Rules' for 'FirewallPolicyNatRuleCollection': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Rules = &output
	}
	return nil
}
