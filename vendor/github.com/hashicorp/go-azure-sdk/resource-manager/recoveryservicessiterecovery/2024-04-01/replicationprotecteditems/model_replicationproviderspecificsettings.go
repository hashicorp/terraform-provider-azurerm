package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationProviderSpecificSettings interface {
	ReplicationProviderSpecificSettings() BaseReplicationProviderSpecificSettingsImpl
}

var _ ReplicationProviderSpecificSettings = BaseReplicationProviderSpecificSettingsImpl{}

type BaseReplicationProviderSpecificSettingsImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseReplicationProviderSpecificSettingsImpl) ReplicationProviderSpecificSettings() BaseReplicationProviderSpecificSettingsImpl {
	return s
}

var _ ReplicationProviderSpecificSettings = RawReplicationProviderSpecificSettingsImpl{}

// RawReplicationProviderSpecificSettingsImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawReplicationProviderSpecificSettingsImpl struct {
	replicationProviderSpecificSettings BaseReplicationProviderSpecificSettingsImpl
	Type                                string
	Values                              map[string]interface{}
}

func (s RawReplicationProviderSpecificSettingsImpl) ReplicationProviderSpecificSettings() BaseReplicationProviderSpecificSettingsImpl {
	return s.replicationProviderSpecificSettings
}

func UnmarshalReplicationProviderSpecificSettingsImplementation(input []byte) (ReplicationProviderSpecificSettings, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ReplicationProviderSpecificSettings into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseReplicationProviderSpecificSettingsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseReplicationProviderSpecificSettingsImpl: %+v", err)
	}

	return RawReplicationProviderSpecificSettingsImpl{
		replicationProviderSpecificSettings: parent,
		Type:                                value,
		Values:                              temp,
	}, nil

}
