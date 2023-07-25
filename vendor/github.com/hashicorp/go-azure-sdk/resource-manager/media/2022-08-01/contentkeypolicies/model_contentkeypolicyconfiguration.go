package contentkeypolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentKeyPolicyConfiguration interface {
}

func unmarshalContentKeyPolicyConfigurationImplementation(input []byte) (ContentKeyPolicyConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicyConfiguration into map[string]interface: %+v", err)
	}

	value, ok := temp["@odata.type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.ContentKeyPolicyClearKeyConfiguration") {
		var out ContentKeyPolicyClearKeyConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContentKeyPolicyClearKeyConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.ContentKeyPolicyFairPlayConfiguration") {
		var out ContentKeyPolicyFairPlayConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContentKeyPolicyFairPlayConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.ContentKeyPolicyPlayReadyConfiguration") {
		var out ContentKeyPolicyPlayReadyConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContentKeyPolicyPlayReadyConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.ContentKeyPolicyUnknownConfiguration") {
		var out ContentKeyPolicyUnknownConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContentKeyPolicyUnknownConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.ContentKeyPolicyWidevineConfiguration") {
		var out ContentKeyPolicyWidevineConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContentKeyPolicyWidevineConfiguration: %+v", err)
		}
		return out, nil
	}

	type RawContentKeyPolicyConfigurationImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawContentKeyPolicyConfigurationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
