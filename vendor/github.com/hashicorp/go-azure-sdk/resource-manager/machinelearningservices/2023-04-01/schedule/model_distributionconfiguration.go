package schedule

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DistributionConfiguration interface {
}

func unmarshalDistributionConfigurationImplementation(input []byte) (DistributionConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DistributionConfiguration into map[string]interface: %+v", err)
	}

	value, ok := temp["distributionType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Mpi") {
		var out Mpi
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Mpi: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PyTorch") {
		var out PyTorch
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PyTorch: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TensorFlow") {
		var out TensorFlow
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TensorFlow: %+v", err)
		}
		return out, nil
	}

	type RawDistributionConfigurationImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawDistributionConfigurationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
