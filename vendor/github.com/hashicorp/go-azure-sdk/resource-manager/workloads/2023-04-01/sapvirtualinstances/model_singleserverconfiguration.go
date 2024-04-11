package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ InfrastructureConfiguration = SingleServerConfiguration{}

type SingleServerConfiguration struct {
	CustomResourceNames         SingleServerCustomResourceNames `json:"customResourceNames"`
	DatabaseType                *SAPDatabaseType                `json:"databaseType,omitempty"`
	DbDiskConfiguration         *DiskConfiguration              `json:"dbDiskConfiguration,omitempty"`
	NetworkConfiguration        *NetworkConfiguration           `json:"networkConfiguration,omitempty"`
	SubnetId                    string                          `json:"subnetId"`
	VirtualMachineConfiguration VirtualMachineConfiguration     `json:"virtualMachineConfiguration"`

	// Fields inherited from InfrastructureConfiguration
	AppResourceGroup string `json:"appResourceGroup"`
}

var _ json.Marshaler = SingleServerConfiguration{}

func (s SingleServerConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper SingleServerConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SingleServerConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SingleServerConfiguration: %+v", err)
	}
	decoded["deploymentType"] = "SingleServer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SingleServerConfiguration: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &SingleServerConfiguration{}

func (s *SingleServerConfiguration) UnmarshalJSON(bytes []byte) error {
	type alias SingleServerConfiguration
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into SingleServerConfiguration: %+v", err)
	}

	s.AppResourceGroup = decoded.AppResourceGroup
	s.DatabaseType = decoded.DatabaseType
	s.DbDiskConfiguration = decoded.DbDiskConfiguration
	s.NetworkConfiguration = decoded.NetworkConfiguration
	s.SubnetId = decoded.SubnetId
	s.VirtualMachineConfiguration = decoded.VirtualMachineConfiguration

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SingleServerConfiguration into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["customResourceNames"]; ok {
		impl, err := unmarshalSingleServerCustomResourceNamesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CustomResourceNames' for 'SingleServerConfiguration': %+v", err)
		}
		s.CustomResourceNames = impl
	}
	return nil
}
