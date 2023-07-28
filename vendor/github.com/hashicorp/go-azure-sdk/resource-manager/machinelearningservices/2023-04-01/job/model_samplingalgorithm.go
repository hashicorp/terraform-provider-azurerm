package job

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SamplingAlgorithm interface {
}

func unmarshalSamplingAlgorithmImplementation(input []byte) (SamplingAlgorithm, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SamplingAlgorithm into map[string]interface: %+v", err)
	}

	value, ok := temp["samplingAlgorithmType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Bayesian") {
		var out BayesianSamplingAlgorithm
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BayesianSamplingAlgorithm: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Grid") {
		var out GridSamplingAlgorithm
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GridSamplingAlgorithm: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Random") {
		var out RandomSamplingAlgorithm
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RandomSamplingAlgorithm: %+v", err)
		}
		return out, nil
	}

	type RawSamplingAlgorithmImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawSamplingAlgorithmImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
