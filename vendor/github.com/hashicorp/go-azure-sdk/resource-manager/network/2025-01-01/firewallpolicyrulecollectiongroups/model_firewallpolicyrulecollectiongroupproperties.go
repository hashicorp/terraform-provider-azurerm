package firewallpolicyrulecollectiongroups

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyRuleCollectionGroupProperties struct {
	Priority          *int64                          `json:"priority,omitempty"`
	ProvisioningState *ProvisioningState              `json:"provisioningState,omitempty"`
	RuleCollections   *[]FirewallPolicyRuleCollection `json:"ruleCollections,omitempty"`
	Size              *string                         `json:"size,omitempty"`
}

var _ json.Unmarshaler = &FirewallPolicyRuleCollectionGroupProperties{}

func (s *FirewallPolicyRuleCollectionGroupProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Priority          *int64             `json:"priority,omitempty"`
		ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
		Size              *string            `json:"size,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Priority = decoded.Priority
	s.ProvisioningState = decoded.ProvisioningState
	s.Size = decoded.Size

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling FirewallPolicyRuleCollectionGroupProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["ruleCollections"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling RuleCollections into list []json.RawMessage: %+v", err)
		}

		output := make([]FirewallPolicyRuleCollection, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalFirewallPolicyRuleCollectionImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'RuleCollections' for 'FirewallPolicyRuleCollectionGroupProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.RuleCollections = &output
	}

	return nil
}
