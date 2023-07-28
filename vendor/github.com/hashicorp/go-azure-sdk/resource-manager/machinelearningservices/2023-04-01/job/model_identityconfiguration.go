package job

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IdentityConfiguration interface {
}

func unmarshalIdentityConfigurationImplementation(input []byte) (IdentityConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling IdentityConfiguration into map[string]interface: %+v", err)
	}

	value, ok := temp["identityType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AMLToken") {
		var out AmlToken
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmlToken: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Managed") {
		var out ManagedIdentity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ManagedIdentity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "UserIdentity") {
		var out UserIdentity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UserIdentity: %+v", err)
		}
		return out, nil
	}

	type RawIdentityConfigurationImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawIdentityConfigurationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
