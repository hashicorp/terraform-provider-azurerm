package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomPersistentDiskProperties interface {
	CustomPersistentDiskProperties() BaseCustomPersistentDiskPropertiesImpl
}

var _ CustomPersistentDiskProperties = BaseCustomPersistentDiskPropertiesImpl{}

type BaseCustomPersistentDiskPropertiesImpl struct {
	EnableSubPath *bool     `json:"enableSubPath,omitempty"`
	MountOptions  *[]string `json:"mountOptions,omitempty"`
	MountPath     string    `json:"mountPath"`
	ReadOnly      *bool     `json:"readOnly,omitempty"`
	Type          Type      `json:"type"`
}

func (s BaseCustomPersistentDiskPropertiesImpl) CustomPersistentDiskProperties() BaseCustomPersistentDiskPropertiesImpl {
	return s
}

var _ CustomPersistentDiskProperties = RawCustomPersistentDiskPropertiesImpl{}

// RawCustomPersistentDiskPropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawCustomPersistentDiskPropertiesImpl struct {
	customPersistentDiskProperties BaseCustomPersistentDiskPropertiesImpl
	Type                           string
	Values                         map[string]interface{}
}

func (s RawCustomPersistentDiskPropertiesImpl) CustomPersistentDiskProperties() BaseCustomPersistentDiskPropertiesImpl {
	return s.customPersistentDiskProperties
}

func (s RawCustomPersistentDiskPropertiesImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalCustomPersistentDiskPropertiesImplementation(input []byte) (CustomPersistentDiskProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomPersistentDiskProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureFileVolume") {
		var out AzureFileVolume
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFileVolume: %+v", err)
		}
		return out, nil
	}

	var parent BaseCustomPersistentDiskPropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseCustomPersistentDiskPropertiesImpl: %+v", err)
	}

	return RawCustomPersistentDiskPropertiesImpl{
		customPersistentDiskProperties: parent,
		Type:                           value,
		Values:                         temp,
	}, nil

}
