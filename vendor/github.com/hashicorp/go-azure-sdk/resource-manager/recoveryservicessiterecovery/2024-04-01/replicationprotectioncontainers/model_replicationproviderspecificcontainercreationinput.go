package replicationprotectioncontainers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationProviderSpecificContainerCreationInput interface {
	ReplicationProviderSpecificContainerCreationInput() BaseReplicationProviderSpecificContainerCreationInputImpl
}

var _ ReplicationProviderSpecificContainerCreationInput = BaseReplicationProviderSpecificContainerCreationInputImpl{}

type BaseReplicationProviderSpecificContainerCreationInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseReplicationProviderSpecificContainerCreationInputImpl) ReplicationProviderSpecificContainerCreationInput() BaseReplicationProviderSpecificContainerCreationInputImpl {
	return s
}

var _ ReplicationProviderSpecificContainerCreationInput = RawReplicationProviderSpecificContainerCreationInputImpl{}

// RawReplicationProviderSpecificContainerCreationInputImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawReplicationProviderSpecificContainerCreationInputImpl struct {
	replicationProviderSpecificContainerCreationInput BaseReplicationProviderSpecificContainerCreationInputImpl
	Type                                              string
	Values                                            map[string]interface{}
}

func (s RawReplicationProviderSpecificContainerCreationInputImpl) ReplicationProviderSpecificContainerCreationInput() BaseReplicationProviderSpecificContainerCreationInputImpl {
	return s.replicationProviderSpecificContainerCreationInput
}

func (s RawReplicationProviderSpecificContainerCreationInputImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalReplicationProviderSpecificContainerCreationInputImplementation(input []byte) (ReplicationProviderSpecificContainerCreationInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ReplicationProviderSpecificContainerCreationInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out A2AContainerCreationInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2AContainerCreationInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "A2ACrossClusterMigration") {
		var out A2ACrossClusterMigrationContainerCreationInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2ACrossClusterMigrationContainerCreationInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VMwareCbt") {
		var out VMwareCbtContainerCreationInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMwareCbtContainerCreationInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseReplicationProviderSpecificContainerCreationInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseReplicationProviderSpecificContainerCreationInputImpl: %+v", err)
	}

	return RawReplicationProviderSpecificContainerCreationInputImpl{
		replicationProviderSpecificContainerCreationInput: parent,
		Type:   value,
		Values: temp,
	}, nil

}
