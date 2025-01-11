package backuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BasePolicyRule interface {
	BasePolicyRule() BaseBasePolicyRuleImpl
}

var _ BasePolicyRule = BaseBasePolicyRuleImpl{}

type BaseBasePolicyRuleImpl struct {
	Name       string `json:"name"`
	ObjectType string `json:"objectType"`
}

func (s BaseBasePolicyRuleImpl) BasePolicyRule() BaseBasePolicyRuleImpl {
	return s
}

var _ BasePolicyRule = RawBasePolicyRuleImpl{}

// RawBasePolicyRuleImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawBasePolicyRuleImpl struct {
	basePolicyRule BaseBasePolicyRuleImpl
	Type           string
	Values         map[string]interface{}
}

func (s RawBasePolicyRuleImpl) BasePolicyRule() BaseBasePolicyRuleImpl {
	return s.basePolicyRule
}

func UnmarshalBasePolicyRuleImplementation(input []byte) (BasePolicyRule, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BasePolicyRule into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureBackupRule") {
		var out AzureBackupRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBackupRule: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureRetentionRule") {
		var out AzureRetentionRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureRetentionRule: %+v", err)
		}
		return out, nil
	}

	var parent BaseBasePolicyRuleImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseBasePolicyRuleImpl: %+v", err)
	}

	return RawBasePolicyRuleImpl{
		basePolicyRule: parent,
		Type:           value,
		Values:         temp,
	}, nil

}
