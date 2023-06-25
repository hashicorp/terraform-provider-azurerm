package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PlannedFailoverProviderSpecificFailoverInput interface {
}

func unmarshalPlannedFailoverProviderSpecificFailoverInputImplementation(input []byte) (PlannedFailoverProviderSpecificFailoverInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling PlannedFailoverProviderSpecificFailoverInput into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "HyperVReplicaAzureFailback") {
		var out HyperVReplicaAzureFailbackProviderInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaAzureFailbackProviderInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplicaAzure") {
		var out HyperVReplicaAzurePlannedFailoverProviderInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaAzurePlannedFailoverProviderInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcmFailback") {
		var out InMageRcmFailbackPlannedFailoverProviderInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmFailbackPlannedFailoverProviderInput: %+v", err)
		}
		return out, nil
	}

	type RawPlannedFailoverProviderSpecificFailoverInputImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawPlannedFailoverProviderSpecificFailoverInputImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
