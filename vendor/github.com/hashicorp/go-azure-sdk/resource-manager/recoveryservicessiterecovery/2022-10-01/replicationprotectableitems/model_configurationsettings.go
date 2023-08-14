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

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawConfigurationSettingsImpl struct {
	Type   string
	Values map[string]interface{}
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

	out := RawConfigurationSettingsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
