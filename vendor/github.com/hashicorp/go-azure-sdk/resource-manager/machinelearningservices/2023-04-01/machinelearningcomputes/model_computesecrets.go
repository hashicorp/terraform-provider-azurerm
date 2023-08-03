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

	type RawComputeSecretsImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawComputeSecretsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
