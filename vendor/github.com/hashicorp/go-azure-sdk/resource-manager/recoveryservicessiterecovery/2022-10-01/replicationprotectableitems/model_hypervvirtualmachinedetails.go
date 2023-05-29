package replicationprotectableitems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConfigurationSettings = HyperVVirtualMachineDetails{}

type HyperVVirtualMachineDetails struct {
	DiskDetails            *[]DiskDetails  `json:"diskDetails,omitempty"`
	Generation             *string         `json:"generation,omitempty"`
	HasFibreChannelAdapter *PresenceStatus `json:"hasFibreChannelAdapter,omitempty"`
	HasPhysicalDisk        *PresenceStatus `json:"hasPhysicalDisk,omitempty"`
	HasSharedVhd           *PresenceStatus `json:"hasSharedVhd,omitempty"`
	HyperVHostId           *string         `json:"hyperVHostId,omitempty"`
	OsDetails              *OSDetails      `json:"osDetails,omitempty"`
	SourceItemId           *string         `json:"sourceItemId,omitempty"`

	// Fields inherited from ConfigurationSettings
}

var _ json.Marshaler = HyperVVirtualMachineDetails{}

func (s HyperVVirtualMachineDetails) MarshalJSON() ([]byte, error) {
	type wrapper HyperVVirtualMachineDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVVirtualMachineDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVVirtualMachineDetails: %+v", err)
	}
	decoded["instanceType"] = "HyperVVirtualMachine"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVVirtualMachineDetails: %+v", err)
	}

	return encoded, nil
}
