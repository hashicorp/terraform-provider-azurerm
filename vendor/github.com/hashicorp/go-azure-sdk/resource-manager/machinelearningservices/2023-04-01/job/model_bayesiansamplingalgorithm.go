package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SamplingAlgorithm = BayesianSamplingAlgorithm{}

type BayesianSamplingAlgorithm struct {

	// Fields inherited from SamplingAlgorithm
}

var _ json.Marshaler = BayesianSamplingAlgorithm{}

func (s BayesianSamplingAlgorithm) MarshalJSON() ([]byte, error) {
	type wrapper BayesianSamplingAlgorithm
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BayesianSamplingAlgorithm: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BayesianSamplingAlgorithm: %+v", err)
	}
	decoded["samplingAlgorithmType"] = "Bayesian"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BayesianSamplingAlgorithm: %+v", err)
	}

	return encoded, nil
}
