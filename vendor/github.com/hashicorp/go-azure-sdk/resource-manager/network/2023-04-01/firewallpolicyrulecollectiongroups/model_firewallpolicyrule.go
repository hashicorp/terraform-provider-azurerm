package firewallpolicyrulecollectiongroups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyRule interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFirewallPolicyRuleImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalFirewallPolicyRuleImplementation(input []byte) (FirewallPolicyRule, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FirewallPolicyRule into map[string]interface: %+v", err)
	}

	value, ok := temp["ruleType"].(string)
	if !ok {
		return nil, nil
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

	out := RawFirewallPolicyRuleImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
