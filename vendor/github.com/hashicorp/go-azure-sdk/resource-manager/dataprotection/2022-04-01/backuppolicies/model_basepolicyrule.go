package backuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BasePolicyRule interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawBasePolicyRuleImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalBasePolicyRuleImplementation(input []byte) (BasePolicyRule, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BasePolicyRule into map[string]interface: %+v", err)
	}

	value, ok := temp["objectType"].(string)
	if !ok {
		return nil, nil
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

	out := RawBasePolicyRuleImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
