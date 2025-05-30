package integrationruntimes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeStatus interface {
	IntegrationRuntimeStatus() BaseIntegrationRuntimeStatusImpl
}

var _ IntegrationRuntimeStatus = BaseIntegrationRuntimeStatusImpl{}

type BaseIntegrationRuntimeStatusImpl struct {
	DataFactoryName *string                  `json:"dataFactoryName,omitempty"`
	State           *IntegrationRuntimeState `json:"state,omitempty"`
	Type            IntegrationRuntimeType   `json:"type"`
}

func (s BaseIntegrationRuntimeStatusImpl) IntegrationRuntimeStatus() BaseIntegrationRuntimeStatusImpl {
	return s
}

var _ IntegrationRuntimeStatus = RawIntegrationRuntimeStatusImpl{}

// RawIntegrationRuntimeStatusImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawIntegrationRuntimeStatusImpl struct {
	integrationRuntimeStatus BaseIntegrationRuntimeStatusImpl
	Type                     string
	Values                   map[string]interface{}
}

func (s RawIntegrationRuntimeStatusImpl) IntegrationRuntimeStatus() BaseIntegrationRuntimeStatusImpl {
	return s.integrationRuntimeStatus
}

func UnmarshalIntegrationRuntimeStatusImplementation(input []byte) (IntegrationRuntimeStatus, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling IntegrationRuntimeStatus into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Managed") {
		var out ManagedIntegrationRuntimeStatus
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ManagedIntegrationRuntimeStatus: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SelfHosted") {
		var out SelfHostedIntegrationRuntimeStatus
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SelfHostedIntegrationRuntimeStatus: %+v", err)
		}
		return out, nil
	}

	var parent BaseIntegrationRuntimeStatusImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseIntegrationRuntimeStatusImpl: %+v", err)
	}

	return RawIntegrationRuntimeStatusImpl{
		integrationRuntimeStatus: parent,
		Type:                     value,
		Values:                   temp,
	}, nil

}
