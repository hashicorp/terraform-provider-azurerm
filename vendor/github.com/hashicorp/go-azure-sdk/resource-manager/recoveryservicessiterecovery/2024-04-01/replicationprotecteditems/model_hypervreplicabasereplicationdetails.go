package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificSettings = HyperVReplicaBaseReplicationDetails{}

type HyperVReplicaBaseReplicationDetails struct {
	InitialReplicationDetails    *InitialReplicationDetails `json:"initialReplicationDetails,omitempty"`
	LastReplicatedTime           *string                    `json:"lastReplicatedTime,omitempty"`
	VMDiskDetails                *[]DiskDetails             `json:"vMDiskDetails,omitempty"`
	VMId                         *string                    `json:"vmId,omitempty"`
	VMNics                       *[]VMNicDetails            `json:"vmNics,omitempty"`
	VMProtectionState            *string                    `json:"vmProtectionState,omitempty"`
	VMProtectionStateDescription *string                    `json:"vmProtectionStateDescription,omitempty"`

	// Fields inherited from ReplicationProviderSpecificSettings

	InstanceType string `json:"instanceType"`
}

func (s HyperVReplicaBaseReplicationDetails) ReplicationProviderSpecificSettings() BaseReplicationProviderSpecificSettingsImpl {
	return BaseReplicationProviderSpecificSettingsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = HyperVReplicaBaseReplicationDetails{}

func (s HyperVReplicaBaseReplicationDetails) MarshalJSON() ([]byte, error) {
	type wrapper HyperVReplicaBaseReplicationDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVReplicaBaseReplicationDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVReplicaBaseReplicationDetails: %+v", err)
	}

	decoded["instanceType"] = "HyperVReplicaBaseReplicationDetails"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVReplicaBaseReplicationDetails: %+v", err)
	}

	return encoded, nil
}
