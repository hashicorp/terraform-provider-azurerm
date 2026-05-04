package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SAPConfiguration = DiscoveryConfiguration{}

type DiscoveryConfiguration struct {
	AppLocation                 *string `json:"appLocation,omitempty"`
	CentralServerVMId           *string `json:"centralServerVmId,omitempty"`
	ManagedRgStorageAccountName *string `json:"managedRgStorageAccountName,omitempty"`

	// Fields inherited from SAPConfiguration

	ConfigurationType SAPConfigurationType `json:"configurationType"`
}

func (s DiscoveryConfiguration) SAPConfiguration() BaseSAPConfigurationImpl {
	return BaseSAPConfigurationImpl{
		ConfigurationType: s.ConfigurationType,
	}
}

var _ json.Marshaler = DiscoveryConfiguration{}

func (s DiscoveryConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper DiscoveryConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DiscoveryConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DiscoveryConfiguration: %+v", err)
	}

	decoded["configurationType"] = "Discovery"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DiscoveryConfiguration: %+v", err)
	}

	return encoded, nil
}
