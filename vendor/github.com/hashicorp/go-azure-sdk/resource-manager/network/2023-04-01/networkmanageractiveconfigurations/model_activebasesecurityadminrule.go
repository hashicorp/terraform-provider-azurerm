package networkmanageractiveconfigurations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActiveBaseSecurityAdminRule interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawActiveBaseSecurityAdminRuleImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalActiveBaseSecurityAdminRuleImplementation(input []byte) (ActiveBaseSecurityAdminRule, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ActiveBaseSecurityAdminRule into map[string]interface: %+v", err)
	}

	value, ok := temp["kind"].(string)
	if !ok {
		return nil, nil
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

	out := RawActiveBaseSecurityAdminRuleImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
