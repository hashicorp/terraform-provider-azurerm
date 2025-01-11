package backupinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BaseResourceProperties interface {
	BaseResourceProperties() BaseBaseResourcePropertiesImpl
}

var _ BaseResourceProperties = BaseBaseResourcePropertiesImpl{}

type BaseBaseResourcePropertiesImpl struct {
	ObjectType ResourcePropertiesObjectType `json:"objectType"`
}

func (s BaseBaseResourcePropertiesImpl) BaseResourceProperties() BaseBaseResourcePropertiesImpl {
	return s
}

var _ BaseResourceProperties = RawBaseResourcePropertiesImpl{}

// RawBaseResourcePropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawBaseResourcePropertiesImpl struct {
	baseResourceProperties BaseBaseResourcePropertiesImpl
	Type                   string
	Values                 map[string]interface{}
}

func (s RawBaseResourcePropertiesImpl) BaseResourceProperties() BaseBaseResourcePropertiesImpl {
	return s.baseResourceProperties
}

func UnmarshalBaseResourcePropertiesImplementation(input []byte) (BaseResourceProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BaseResourceProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DefaultResourceProperties") {
		var out DefaultResourceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DefaultResourceProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseBaseResourcePropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseBaseResourcePropertiesImpl: %+v", err)
	}

	return RawBaseResourcePropertiesImpl{
		baseResourceProperties: parent,
		Type:                   value,
		Values:                 temp,
	}, nil

}
