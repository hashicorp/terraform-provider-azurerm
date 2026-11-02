package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OSConfiguration interface {
	OSConfiguration() BaseOSConfigurationImpl
}

var _ OSConfiguration = BaseOSConfigurationImpl{}

type BaseOSConfigurationImpl struct {
	OsType OSType `json:"osType"`
}

func (s BaseOSConfigurationImpl) OSConfiguration() BaseOSConfigurationImpl {
	return s
}

var _ OSConfiguration = RawOSConfigurationImpl{}

// RawOSConfigurationImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawOSConfigurationImpl struct {
	oSConfiguration BaseOSConfigurationImpl
	Type            string
	Values          map[string]interface{}
}

func (s RawOSConfigurationImpl) OSConfiguration() BaseOSConfigurationImpl {
	return s.oSConfiguration
}

func (s RawOSConfigurationImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalOSConfigurationImplementation(input []byte) (OSConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling OSConfiguration into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["osType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Linux") {
		var out LinuxConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LinuxConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Windows") {
		var out WindowsConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WindowsConfiguration: %+v", err)
		}
		return out, nil
	}

	var parent BaseOSConfigurationImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseOSConfigurationImpl: %+v", err)
	}

	return RawOSConfigurationImpl{
		oSConfiguration: parent,
		Type:            value,
		Values:          temp,
	}, nil

}
