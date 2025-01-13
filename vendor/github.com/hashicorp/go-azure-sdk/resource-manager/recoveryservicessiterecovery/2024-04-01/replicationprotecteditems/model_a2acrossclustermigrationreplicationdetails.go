package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificSettings = A2ACrossClusterMigrationReplicationDetails{}

type A2ACrossClusterMigrationReplicationDetails struct {
	FabricObjectId               *string `json:"fabricObjectId,omitempty"`
	LifecycleId                  *string `json:"lifecycleId,omitempty"`
	OsType                       *string `json:"osType,omitempty"`
	PrimaryFabricLocation        *string `json:"primaryFabricLocation,omitempty"`
	VMProtectionState            *string `json:"vmProtectionState,omitempty"`
	VMProtectionStateDescription *string `json:"vmProtectionStateDescription,omitempty"`

	// Fields inherited from ReplicationProviderSpecificSettings

	InstanceType string `json:"instanceType"`
}

func (s A2ACrossClusterMigrationReplicationDetails) ReplicationProviderSpecificSettings() BaseReplicationProviderSpecificSettingsImpl {
	return BaseReplicationProviderSpecificSettingsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = A2ACrossClusterMigrationReplicationDetails{}

func (s A2ACrossClusterMigrationReplicationDetails) MarshalJSON() ([]byte, error) {
	type wrapper A2ACrossClusterMigrationReplicationDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2ACrossClusterMigrationReplicationDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2ACrossClusterMigrationReplicationDetails: %+v", err)
	}

	decoded["instanceType"] = "A2ACrossClusterMigration"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2ACrossClusterMigrationReplicationDetails: %+v", err)
	}

	return encoded, nil
}
