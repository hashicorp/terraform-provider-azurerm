package replicationprotectableitems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationSettings interface {
}

func unmarshalConfigurationSettingsImplementation(input []byte) (ConfigurationSettings, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ConfigurationSettings into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "HyperVVirtualMachine") {
		var out HyperVVirtualMachineDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVVirtualMachineDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ReplicationGroupDetails") {
		var out ReplicationGroupDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ReplicationGroupDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VMwareVirtualMachine") {
		var out VMwareVirtualMachineDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMwareVirtualMachineDetails: %+v", err)
		}
		return out, nil
	}

	type RawConfigurationSettingsImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawConfigurationSettingsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
