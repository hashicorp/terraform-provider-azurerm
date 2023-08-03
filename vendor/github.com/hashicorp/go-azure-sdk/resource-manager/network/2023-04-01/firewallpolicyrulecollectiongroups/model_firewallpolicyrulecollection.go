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

	type RawFirewallPolicyRuleCollectionImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawFirewallPolicyRuleCollectionImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
