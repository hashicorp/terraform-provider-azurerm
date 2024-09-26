package links

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TargetServiceBase interface {
	TargetServiceBase() BaseTargetServiceBaseImpl
}

var _ TargetServiceBase = BaseTargetServiceBaseImpl{}

type BaseTargetServiceBaseImpl struct {
	Type TargetServiceType `json:"type"`
}

func (s BaseTargetServiceBaseImpl) TargetServiceBase() BaseTargetServiceBaseImpl {
	return s
}

var _ TargetServiceBase = RawTargetServiceBaseImpl{}

// RawTargetServiceBaseImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawTargetServiceBaseImpl struct {
	targetServiceBase BaseTargetServiceBaseImpl
	Type              string
	Values            map[string]interface{}
}

func (s RawTargetServiceBaseImpl) TargetServiceBase() BaseTargetServiceBaseImpl {
	return s.targetServiceBase
}

func UnmarshalTargetServiceBaseImplementation(input []byte) (TargetServiceBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TargetServiceBase into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureResource") {
		var out AzureResource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureResource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConfluentBootstrapServer") {
		var out ConfluentBootstrapServer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConfluentBootstrapServer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConfluentSchemaRegistry") {
		var out ConfluentSchemaRegistry
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConfluentSchemaRegistry: %+v", err)
		}
		return out, nil
	}

	var parent BaseTargetServiceBaseImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseTargetServiceBaseImpl: %+v", err)
	}

	return RawTargetServiceBaseImpl{
		targetServiceBase: parent,
		Type:              value,
		Values:            temp,
	}, nil

}
