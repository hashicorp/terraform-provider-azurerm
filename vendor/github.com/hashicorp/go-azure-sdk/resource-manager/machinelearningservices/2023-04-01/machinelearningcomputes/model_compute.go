package machinelearningcomputes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Compute interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawComputeImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalComputeImplementation(input []byte) (Compute, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Compute into map[string]interface: %+v", err)
	}

	value, ok := temp["computeType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AKS") {
		var out AKS
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AKS: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmlCompute") {
		var out AmlCompute
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmlCompute: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ComputeInstance") {
		var out ComputeInstance
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ComputeInstance: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DataFactory") {
		var out DataFactory
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DataFactory: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DataLakeAnalytics") {
		var out DataLakeAnalytics
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DataLakeAnalytics: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Databricks") {
		var out Databricks
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Databricks: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HDInsight") {
		var out HDInsight
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HDInsight: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Kubernetes") {
		var out Kubernetes
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Kubernetes: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SynapseSpark") {
		var out SynapseSpark
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SynapseSpark: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VirtualMachine") {
		var out VirtualMachine
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VirtualMachine: %+v", err)
		}
		return out, nil
	}

	out := RawComputeImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
