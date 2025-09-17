package firewallpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyRule interface {
	FirewallPolicyRule() BaseFirewallPolicyRuleImpl
}

var _ FirewallPolicyRule = BaseFirewallPolicyRuleImpl{}

type BaseFirewallPolicyRuleImpl struct {
	Description *string                `json:"description,omitempty"`
	Name        *string                `json:"name,omitempty"`
	RuleType    FirewallPolicyRuleType `json:"ruleType"`
}

func (s BaseFirewallPolicyRuleImpl) FirewallPolicyRule() BaseFirewallPolicyRuleImpl {
	return s
}

var _ FirewallPolicyRule = RawFirewallPolicyRuleImpl{}

// RawFirewallPolicyRuleImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFirewallPolicyRuleImpl struct {
	firewallPolicyRule BaseFirewallPolicyRuleImpl
	Type               string
	Values             map[string]interface{}
}

func (s RawFirewallPolicyRuleImpl) FirewallPolicyRule() BaseFirewallPolicyRuleImpl {
	return s.firewallPolicyRule
}

func UnmarshalFirewallPolicyRuleImplementation(input []byte) (FirewallPolicyRule, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FirewallPolicyRule into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["ruleType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "ApplicationRule") {
		var out ApplicationRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ApplicationRule: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NatRule") {
		var out NatRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NatRule: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NetworkRule") {
		var out NetworkRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NetworkRule: %+v", err)
		}
		return out, nil
	}

	var parent BaseFirewallPolicyRuleImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFirewallPolicyRuleImpl: %+v", err)
	}

	return RawFirewallPolicyRuleImpl{
		firewallPolicyRule: parent,
		Type:               value,
		Values:             temp,
	}, nil

}
