package firewallpolicyrulecollectiongroups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyRuleCollection interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFirewallPolicyRuleCollectionImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalFirewallPolicyRuleCollectionImplementation(input []byte) (FirewallPolicyRuleCollection, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FirewallPolicyRuleCollection into map[string]interface: %+v", err)
	}

	value, ok := temp["ruleCollectionType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "FirewallPolicyFilterRuleCollection") {
		var out FirewallPolicyFilterRuleCollection
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FirewallPolicyFilterRuleCollection: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FirewallPolicyNatRuleCollection") {
		var out FirewallPolicyNatRuleCollection
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FirewallPolicyNatRuleCollection: %+v", err)
		}
		return out, nil
	}

	out := RawFirewallPolicyRuleCollectionImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
