package securitypolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityPolicyPropertiesParameters interface {
	SecurityPolicyPropertiesParameters() BaseSecurityPolicyPropertiesParametersImpl
}

var _ SecurityPolicyPropertiesParameters = BaseSecurityPolicyPropertiesParametersImpl{}

type BaseSecurityPolicyPropertiesParametersImpl struct {
	Type SecurityPolicyType `json:"type"`
}

func (s BaseSecurityPolicyPropertiesParametersImpl) SecurityPolicyPropertiesParameters() BaseSecurityPolicyPropertiesParametersImpl {
	return s
}

var _ SecurityPolicyPropertiesParameters = RawSecurityPolicyPropertiesParametersImpl{}

// RawSecurityPolicyPropertiesParametersImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawSecurityPolicyPropertiesParametersImpl struct {
	securityPolicyPropertiesParameters BaseSecurityPolicyPropertiesParametersImpl
	Type                               string
	Values                             map[string]interface{}
}

func (s RawSecurityPolicyPropertiesParametersImpl) SecurityPolicyPropertiesParameters() BaseSecurityPolicyPropertiesParametersImpl {
	return s.securityPolicyPropertiesParameters
}

func UnmarshalSecurityPolicyPropertiesParametersImplementation(input []byte) (SecurityPolicyPropertiesParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SecurityPolicyPropertiesParameters into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "WebApplicationFirewall") {
		var out SecurityPolicyWebApplicationFirewallParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SecurityPolicyWebApplicationFirewallParameters: %+v", err)
		}
		return out, nil
	}

	var parent BaseSecurityPolicyPropertiesParametersImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSecurityPolicyPropertiesParametersImpl: %+v", err)
	}

	return RawSecurityPolicyPropertiesParametersImpl{
		securityPolicyPropertiesParameters: parent,
		Type:                               value,
		Values:                             temp,
	}, nil

}
