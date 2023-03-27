package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareConfiguration interface {
}

func unmarshalSoftwareConfigurationImplementation(input []byte) (SoftwareConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SoftwareConfiguration into map[string]interface: %+v", err)
	}

	value, ok := temp["softwareInstallationType"].(string)
	if !ok {
		return nil, nil
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

	type RawSoftwareConfigurationImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawSoftwareConfigurationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
