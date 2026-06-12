package networkmanagereffectivesecurityadminrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EffectiveBaseSecurityAdminRule interface {
	EffectiveBaseSecurityAdminRule() BaseEffectiveBaseSecurityAdminRuleImpl
}

var _ EffectiveBaseSecurityAdminRule = BaseEffectiveBaseSecurityAdminRuleImpl{}

type BaseEffectiveBaseSecurityAdminRuleImpl struct {
	ConfigurationDescription      *string                            `json:"configurationDescription,omitempty"`
	Id                            *string                            `json:"id,omitempty"`
	Kind                          EffectiveAdminRuleKind             `json:"kind"`
	RuleCollectionAppliesToGroups *[]NetworkManagerSecurityGroupItem `json:"ruleCollectionAppliesToGroups,omitempty"`
	RuleCollectionDescription     *string                            `json:"ruleCollectionDescription,omitempty"`
	RuleGroups                    *[]ConfigurationGroup              `json:"ruleGroups,omitempty"`
}

func (s BaseEffectiveBaseSecurityAdminRuleImpl) EffectiveBaseSecurityAdminRule() BaseEffectiveBaseSecurityAdminRuleImpl {
	return s
}

var _ EffectiveBaseSecurityAdminRule = RawEffectiveBaseSecurityAdminRuleImpl{}

// RawEffectiveBaseSecurityAdminRuleImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawEffectiveBaseSecurityAdminRuleImpl struct {
	effectiveBaseSecurityAdminRule BaseEffectiveBaseSecurityAdminRuleImpl
	Type                           string
	Values                         map[string]interface{}
}

func (s RawEffectiveBaseSecurityAdminRuleImpl) EffectiveBaseSecurityAdminRule() BaseEffectiveBaseSecurityAdminRuleImpl {
	return s.effectiveBaseSecurityAdminRule
}

func UnmarshalEffectiveBaseSecurityAdminRuleImplementation(input []byte) (EffectiveBaseSecurityAdminRule, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EffectiveBaseSecurityAdminRule into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Default") {
		var out EffectiveDefaultSecurityAdminRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EffectiveDefaultSecurityAdminRule: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Custom") {
		var out EffectiveSecurityAdminRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EffectiveSecurityAdminRule: %+v", err)
		}
		return out, nil
	}

	var parent BaseEffectiveBaseSecurityAdminRuleImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseEffectiveBaseSecurityAdminRuleImpl: %+v", err)
	}

	return RawEffectiveBaseSecurityAdminRuleImpl{
		effectiveBaseSecurityAdminRule: parent,
		Type:                           value,
		Values:                         temp,
	}, nil

}
