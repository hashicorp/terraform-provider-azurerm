package machinelearningcomputes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeSecrets interface {
	ComputeSecrets() BaseComputeSecretsImpl
}

var _ ComputeSecrets = BaseComputeSecretsImpl{}

type BaseComputeSecretsImpl struct {
	ComputeType ComputeType `json:"computeType"`
}

func (s BaseComputeSecretsImpl) ComputeSecrets() BaseComputeSecretsImpl {
	return s
}

var _ ComputeSecrets = RawComputeSecretsImpl{}

// RawComputeSecretsImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawComputeSecretsImpl struct {
	computeSecrets BaseComputeSecretsImpl
	Type           string
	Values         map[string]interface{}
}

func (s RawComputeSecretsImpl) ComputeSecrets() BaseComputeSecretsImpl {
	return s.computeSecrets
}

func (s RawComputeSecretsImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalComputeSecretsImplementation(input []byte) (ComputeSecrets, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ComputeSecrets into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["computeType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AKS") {
		var out AksComputeSecrets
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AksComputeSecrets: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Databricks") {
		var out DatabricksComputeSecrets
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DatabricksComputeSecrets: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VirtualMachine") {
		var out VirtualMachineSecrets
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VirtualMachineSecrets: %+v", err)
		}
		return out, nil
	}

	var parent BaseComputeSecretsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseComputeSecretsImpl: %+v", err)
	}

	return RawComputeSecretsImpl{
		computeSecrets: parent,
		Type:           value,
		Values:         temp,
	}, nil

}
