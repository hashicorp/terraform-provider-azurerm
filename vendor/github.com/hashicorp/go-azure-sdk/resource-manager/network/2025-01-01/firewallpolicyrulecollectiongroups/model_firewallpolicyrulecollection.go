package firewallpolicyrulecollectiongroups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyRuleCollection interface {
	FirewallPolicyRuleCollection() BaseFirewallPolicyRuleCollectionImpl
}

var _ FirewallPolicyRuleCollection = BaseFirewallPolicyRuleCollectionImpl{}

type BaseFirewallPolicyRuleCollectionImpl struct {
	Name               *string                          `json:"name,omitempty"`
	Priority           *int64                           `json:"priority,omitempty"`
	RuleCollectionType FirewallPolicyRuleCollectionType `json:"ruleCollectionType"`
}

func (s BaseFirewallPolicyRuleCollectionImpl) FirewallPolicyRuleCollection() BaseFirewallPolicyRuleCollectionImpl {
	return s
}

var _ FirewallPolicyRuleCollection = RawFirewallPolicyRuleCollectionImpl{}

// RawFirewallPolicyRuleCollectionImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFirewallPolicyRuleCollectionImpl struct {
	firewallPolicyRuleCollection BaseFirewallPolicyRuleCollectionImpl
	Type                         string
	Values                       map[string]interface{}
}

func (s RawFirewallPolicyRuleCollectionImpl) FirewallPolicyRuleCollection() BaseFirewallPolicyRuleCollectionImpl {
	return s.firewallPolicyRuleCollection
}

func UnmarshalFirewallPolicyRuleCollectionImplementation(input []byte) (FirewallPolicyRuleCollection, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FirewallPolicyRuleCollection into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["ruleCollectionType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseFirewallPolicyRuleCollectionImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFirewallPolicyRuleCollectionImpl: %+v", err)
	}

	return RawFirewallPolicyRuleCollectionImpl{
		firewallPolicyRuleCollection: parent,
		Type:                         value,
		Values:                       temp,
	}, nil

}
