package actionrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionRuleProperties interface {
}

// RawActionRulePropertiesImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawActionRulePropertiesImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalActionRulePropertiesImplementation(input []byte) (ActionRuleProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ActionRuleProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "ActionGroup") {
		var out ActionGroup
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ActionGroup: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Diagnostics") {
		var out Diagnostics
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Diagnostics: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Suppression") {
		var out Suppression
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Suppression: %+v", err)
		}
		return out, nil
	}

	out := RawActionRulePropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
