package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificSettings = HyperVReplicaBlueReplicationDetails{}

type HyperVReplicaBlueReplicationDetails struct {
	InitialReplicationDetails    *InitialReplicationDetails `json:"initialReplicationDetails,omitempty"`
	LastReplicatedTime           *string                    `json:"lastReplicatedTime,omitempty"`
	VMDiskDetails                *[]DiskDetails             `json:"vMDiskDetails,omitempty"`
	VMId                         *string                    `json:"vmId,omitempty"`
	VMNics                       *[]VMNicDetails            `json:"vmNics,omitempty"`
	VMProtectionState            *string                    `json:"vmProtectionState,omitempty"`
	VMProtectionStateDescription *string                    `json:"vmProtectionStateDescription,omitempty"`

	// Fields inherited from ReplicationProviderSpecificSettings
}

var _ json.Marshaler = HyperVReplicaBlueReplicationDetails{}

func (s HyperVReplicaBlueReplicationDetails) MarshalJSON() ([]byte, error) {
	type wrapper HyperVReplicaBlueReplicationDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVReplicaBlueReplicationDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVReplicaBlueReplicationDetails: %+v", err)
	}
	decoded["instanceType"] = "HyperVReplica2012R2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVReplicaBlueReplicationDetails: %+v", err)
	}

	return encoded, nil
}
