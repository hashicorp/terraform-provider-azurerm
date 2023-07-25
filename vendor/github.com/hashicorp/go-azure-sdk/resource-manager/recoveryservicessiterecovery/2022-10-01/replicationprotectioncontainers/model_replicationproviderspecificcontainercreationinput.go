package replicationprotectioncontainers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationProviderSpecificContainerCreationInput interface {
}

func unmarshalReplicationProviderSpecificContainerCreationInputImplementation(input []byte) (ReplicationProviderSpecificContainerCreationInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ReplicationProviderSpecificContainerCreationInput into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
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

	type RawReplicationProviderSpecificContainerCreationInputImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawReplicationProviderSpecificContainerCreationInputImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
