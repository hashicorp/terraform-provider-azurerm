package machinelearningcomputes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeSecrets interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawComputeSecretsImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalComputeSecretsImplementation(input []byte) (ComputeSecrets, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ComputeSecrets into map[string]interface: %+v", err)
	}

	value, ok := temp["computeType"].(string)
	if !ok {
		return nil, nil
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

	out := RawComputeSecretsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
