package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PlannedFailoverProviderSpecificFailoverInput interface {
	PlannedFailoverProviderSpecificFailoverInput() BasePlannedFailoverProviderSpecificFailoverInputImpl
}

var _ PlannedFailoverProviderSpecificFailoverInput = BasePlannedFailoverProviderSpecificFailoverInputImpl{}

type BasePlannedFailoverProviderSpecificFailoverInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BasePlannedFailoverProviderSpecificFailoverInputImpl) PlannedFailoverProviderSpecificFailoverInput() BasePlannedFailoverProviderSpecificFailoverInputImpl {
	return s
}

var _ PlannedFailoverProviderSpecificFailoverInput = RawPlannedFailoverProviderSpecificFailoverInputImpl{}

// RawPlannedFailoverProviderSpecificFailoverInputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawPlannedFailoverProviderSpecificFailoverInputImpl struct {
	plannedFailoverProviderSpecificFailoverInput BasePlannedFailoverProviderSpecificFailoverInputImpl
	Type                                         string
	Values                                       map[string]interface{}
}

func (s RawPlannedFailoverProviderSpecificFailoverInputImpl) PlannedFailoverProviderSpecificFailoverInput() BasePlannedFailoverProviderSpecificFailoverInputImpl {
	return s.plannedFailoverProviderSpecificFailoverInput
}

func UnmarshalPlannedFailoverProviderSpecificFailoverInputImplementation(input []byte) (PlannedFailoverProviderSpecificFailoverInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling PlannedFailoverProviderSpecificFailoverInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BasePlannedFailoverProviderSpecificFailoverInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BasePlannedFailoverProviderSpecificFailoverInputImpl: %+v", err)
	}

	return RawPlannedFailoverProviderSpecificFailoverInputImpl{
		plannedFailoverProviderSpecificFailoverInput: parent,
		Type:   value,
		Values: temp,
	}, nil

}
