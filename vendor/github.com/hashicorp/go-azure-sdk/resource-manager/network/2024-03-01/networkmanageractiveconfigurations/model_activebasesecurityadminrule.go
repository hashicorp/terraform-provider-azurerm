package networkmanageractiveconfigurations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActiveBaseSecurityAdminRule interface {
	ActiveBaseSecurityAdminRule() BaseActiveBaseSecurityAdminRuleImpl
}

var _ ActiveBaseSecurityAdminRule = BaseActiveBaseSecurityAdminRuleImpl{}

type BaseActiveBaseSecurityAdminRuleImpl struct {
	CommitTime                    *string                            `json:"commitTime,omitempty"`
	ConfigurationDescription      *string                            `json:"configurationDescription,omitempty"`
	Id                            *string                            `json:"id,omitempty"`
	Kind                          EffectiveAdminRuleKind             `json:"kind"`
	Region                        *string                            `json:"region,omitempty"`
	RuleCollectionAppliesToGroups *[]NetworkManagerSecurityGroupItem `json:"ruleCollectionAppliesToGroups,omitempty"`
	RuleCollectionDescription     *string                            `json:"ruleCollectionDescription,omitempty"`
	RuleGroups                    *[]ConfigurationGroup              `json:"ruleGroups,omitempty"`
}

func (s BaseActiveBaseSecurityAdminRuleImpl) ActiveBaseSecurityAdminRule() BaseActiveBaseSecurityAdminRuleImpl {
	return s
}

var _ ActiveBaseSecurityAdminRule = RawActiveBaseSecurityAdminRuleImpl{}

// RawActiveBaseSecurityAdminRuleImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawActiveBaseSecurityAdminRuleImpl struct {
	activeBaseSecurityAdminRule BaseActiveBaseSecurityAdminRuleImpl
	Type                        string
	Values                      map[string]interface{}
}

func (s RawActiveBaseSecurityAdminRuleImpl) ActiveBaseSecurityAdminRule() BaseActiveBaseSecurityAdminRuleImpl {
	return s.activeBaseSecurityAdminRule
}

func UnmarshalActiveBaseSecurityAdminRuleImplementation(input []byte) (ActiveBaseSecurityAdminRule, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ActiveBaseSecurityAdminRule into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Default") {
		var out ActiveDefaultSecurityAdminRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ActiveDefaultSecurityAdminRule: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Custom") {
		var out ActiveSecurityAdminRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ActiveSecurityAdminRule: %+v", err)
		}
		return out, nil
	}

	var parent BaseActiveBaseSecurityAdminRuleImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseActiveBaseSecurityAdminRuleImpl: %+v", err)
	}

	return RawActiveBaseSecurityAdminRuleImpl{
		activeBaseSecurityAdminRule: parent,
		Type:                        value,
		Values:                      temp,
	}, nil

}
