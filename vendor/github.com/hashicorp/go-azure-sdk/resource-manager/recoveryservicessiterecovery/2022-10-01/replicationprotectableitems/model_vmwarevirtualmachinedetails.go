package replicationprotectableitems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConfigurationSettings = VMwareVirtualMachineDetails{}

type VMwareVirtualMachineDetails struct {
	AgentGeneratedId        *string              `json:"agentGeneratedId,omitempty"`
	AgentInstalled          *string              `json:"agentInstalled,omitempty"`
	AgentVersion            *string              `json:"agentVersion,omitempty"`
	DiscoveryType           *string              `json:"discoveryType,omitempty"`
	DiskDetails             *[]InMageDiskDetails `json:"diskDetails,omitempty"`
	IPAddress               *string              `json:"ipAddress,omitempty"`
	OsType                  *string              `json:"osType,omitempty"`
	PoweredOn               *string              `json:"poweredOn,omitempty"`
	VCenterInfrastructureId *string              `json:"vCenterInfrastructureId,omitempty"`
	ValidationErrors        *[]HealthError       `json:"validationErrors,omitempty"`

	// Fields inherited from ConfigurationSettings
}

var _ json.Marshaler = VMwareVirtualMachineDetails{}

func (s VMwareVirtualMachineDetails) MarshalJSON() ([]byte, error) {
	type wrapper VMwareVirtualMachineDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMwareVirtualMachineDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMwareVirtualMachineDetails: %+v", err)
	}
	decoded["instanceType"] = "VMwareVirtualMachine"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMwareVirtualMachineDetails: %+v", err)
	}

	return encoded, nil
}
