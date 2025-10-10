package integrationruntimes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntime interface {
	IntegrationRuntime() BaseIntegrationRuntimeImpl
}

var _ IntegrationRuntime = BaseIntegrationRuntimeImpl{}

type BaseIntegrationRuntimeImpl struct {
	Description *string                `json:"description,omitempty"`
	Type        IntegrationRuntimeType `json:"type"`
}

func (s BaseIntegrationRuntimeImpl) IntegrationRuntime() BaseIntegrationRuntimeImpl {
	return s
}

var _ IntegrationRuntime = RawIntegrationRuntimeImpl{}

// RawIntegrationRuntimeImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawIntegrationRuntimeImpl struct {
	integrationRuntime BaseIntegrationRuntimeImpl
	Type               string
	Values             map[string]interface{}
}

func (s RawIntegrationRuntimeImpl) IntegrationRuntime() BaseIntegrationRuntimeImpl {
	return s.integrationRuntime
}

func UnmarshalIntegrationRuntimeImplementation(input []byte) (IntegrationRuntime, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling IntegrationRuntime into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Managed") {
		var out ManagedIntegrationRuntime
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ManagedIntegrationRuntime: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SelfHosted") {
		var out SelfHostedIntegrationRuntime
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SelfHostedIntegrationRuntime: %+v", err)
		}
		return out, nil
	}

	var parent BaseIntegrationRuntimeImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseIntegrationRuntimeImpl: %+v", err)
	}

	return RawIntegrationRuntimeImpl{
		integrationRuntime: parent,
		Type:               value,
		Values:             temp,
	}, nil

}
