package firewallpolicyrulecollectiongroups

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FirewallPolicyRuleCollection = FirewallPolicyFilterRuleCollection{}

type FirewallPolicyFilterRuleCollection struct {
	Action *FirewallPolicyFilterRuleCollectionAction `json:"action,omitempty"`
	Rules  *[]FirewallPolicyRule                     `json:"rules,omitempty"`

	// Fields inherited from FirewallPolicyRuleCollection

	Name               *string                          `json:"name,omitempty"`
	Priority           *int64                           `json:"priority,omitempty"`
	RuleCollectionType FirewallPolicyRuleCollectionType `json:"ruleCollectionType"`
}

func (s FirewallPolicyFilterRuleCollection) FirewallPolicyRuleCollection() BaseFirewallPolicyRuleCollectionImpl {
	return BaseFirewallPolicyRuleCollectionImpl{
		Name:               s.Name,
		Priority:           s.Priority,
		RuleCollectionType: s.RuleCollectionType,
	}
}

var _ json.Marshaler = FirewallPolicyFilterRuleCollection{}

func (s FirewallPolicyFilterRuleCollection) MarshalJSON() ([]byte, error) {
	type wrapper FirewallPolicyFilterRuleCollection
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FirewallPolicyFilterRuleCollection: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FirewallPolicyFilterRuleCollection: %+v", err)
	}

	decoded["ruleCollectionType"] = "FirewallPolicyFilterRuleCollection"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FirewallPolicyFilterRuleCollection: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &FirewallPolicyFilterRuleCollection{}

func (s *FirewallPolicyFilterRuleCollection) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Action             *FirewallPolicyFilterRuleCollectionAction `json:"action,omitempty"`
		Name               *string                                   `json:"name,omitempty"`
		Priority           *int64                                    `json:"priority,omitempty"`
		RuleCollectionType FirewallPolicyRuleCollectionType          `json:"ruleCollectionType"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Action = decoded.Action
	s.Name = decoded.Name
	s.Priority = decoded.Priority
	s.RuleCollectionType = decoded.RuleCollectionType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling FirewallPolicyFilterRuleCollection into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["rules"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Rules into list []json.RawMessage: %+v", err)
		}

		output := make([]FirewallPolicyRule, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalFirewallPolicyRuleImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Rules' for 'FirewallPolicyFilterRuleCollection': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Rules = &output
	}

	return nil
}
