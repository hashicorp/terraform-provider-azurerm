package adminrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BaseAdminRule interface {
}

func unmarshalBaseAdminRuleImplementation(input []byte) (BaseAdminRule, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BaseAdminRule into map[string]interface: %+v", err)
	}

	value, ok := temp["kind"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Custom") {
		var out AdminRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AdminRule: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Default") {
		var out DefaultAdminRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DefaultAdminRule: %+v", err)
		}
		return out, nil
	}

	type RawBaseAdminRuleImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawBaseAdminRuleImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
