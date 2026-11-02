package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareConfiguration interface {
	SoftwareConfiguration() BaseSoftwareConfigurationImpl
}

var _ SoftwareConfiguration = BaseSoftwareConfigurationImpl{}

type BaseSoftwareConfigurationImpl struct {
	SoftwareInstallationType SAPSoftwareInstallationType `json:"softwareInstallationType"`
}

func (s BaseSoftwareConfigurationImpl) SoftwareConfiguration() BaseSoftwareConfigurationImpl {
	return s
}

var _ SoftwareConfiguration = RawSoftwareConfigurationImpl{}

// RawSoftwareConfigurationImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawSoftwareConfigurationImpl struct {
	softwareConfiguration BaseSoftwareConfigurationImpl
	Type                  string
	Values                map[string]interface{}
}

func (s RawSoftwareConfigurationImpl) SoftwareConfiguration() BaseSoftwareConfigurationImpl {
	return s.softwareConfiguration
}

func (s RawSoftwareConfigurationImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalSoftwareConfigurationImplementation(input []byte) (SoftwareConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SoftwareConfiguration into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["softwareInstallationType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "External") {
		var out ExternalInstallationSoftwareConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ExternalInstallationSoftwareConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SAPInstallWithoutOSConfig") {
		var out SAPInstallWithoutOSConfigSoftwareConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SAPInstallWithoutOSConfigSoftwareConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServiceInitiated") {
		var out ServiceInitiatedSoftwareConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServiceInitiatedSoftwareConfiguration: %+v", err)
		}
		return out, nil
	}

	var parent BaseSoftwareConfigurationImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSoftwareConfigurationImpl: %+v", err)
	}

	return RawSoftwareConfigurationImpl{
		softwareConfiguration: parent,
		Type:                  value,
		Values:                temp,
	}, nil

}
