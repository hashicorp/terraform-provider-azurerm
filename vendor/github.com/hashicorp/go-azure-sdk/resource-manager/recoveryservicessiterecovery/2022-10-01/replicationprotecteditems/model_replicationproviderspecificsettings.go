package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationProviderSpecificSettings interface {
}

func unmarshalReplicationProviderSpecificSettingsImplementation(input []byte) (ReplicationProviderSpecificSettings, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ReplicationProviderSpecificSettings into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "A2ACrossClusterMigration") {
		var out A2ACrossClusterMigrationReplicationDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2ACrossClusterMigrationReplicationDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "A2A") {
		var out A2AReplicationDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2AReplicationDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplicaAzure") {
		var out HyperVReplicaAzureReplicationDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaAzureReplicationDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplicaBaseReplicationDetails") {
		var out HyperVReplicaBaseReplicationDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaBaseReplicationDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplica2012R2") {
		var out HyperVReplicaBlueReplicationDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaBlueReplicationDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplica2012") {
		var out HyperVReplicaReplicationDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaReplicationDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageAzureV2") {
		var out InMageAzureV2ReplicationDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageAzureV2ReplicationDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcmFailback") {
		var out InMageRcmFailbackReplicationDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmFailbackReplicationDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmReplicationDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmReplicationDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMage") {
		var out InMageReplicationDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageReplicationDetails: %+v", err)
		}
		return out, nil
	}

	type RawReplicationProviderSpecificSettingsImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawReplicationProviderSpecificSettingsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
