package machinelearningcomputes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Compute interface {
	Compute() BaseComputeImpl
}

var _ Compute = BaseComputeImpl{}

type BaseComputeImpl struct {
	ComputeLocation    *string            `json:"computeLocation,omitempty"`
	ComputeType        ComputeType        `json:"computeType"`
	CreatedOn          *string            `json:"createdOn,omitempty"`
	Description        *string            `json:"description,omitempty"`
	DisableLocalAuth   *bool              `json:"disableLocalAuth,omitempty"`
	IsAttachedCompute  *bool              `json:"isAttachedCompute,omitempty"`
	ModifiedOn         *string            `json:"modifiedOn,omitempty"`
	ProvisioningErrors *[]ErrorResponse   `json:"provisioningErrors,omitempty"`
	ProvisioningState  *ProvisioningState `json:"provisioningState,omitempty"`
	ResourceId         *string            `json:"resourceId,omitempty"`
}

func (s BaseComputeImpl) Compute() BaseComputeImpl {
	return s
}

var _ Compute = RawComputeImpl{}

// RawComputeImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawComputeImpl struct {
	compute BaseComputeImpl
	Type    string
	Values  map[string]interface{}
}

func (s RawComputeImpl) Compute() BaseComputeImpl {
	return s.compute
}

func UnmarshalComputeImplementation(input []byte) (Compute, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Compute into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["computeType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseComputeImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseComputeImpl: %+v", err)
	}

	return RawComputeImpl{
		compute: parent,
		Type:    value,
		Values:  temp,
	}, nil

}
