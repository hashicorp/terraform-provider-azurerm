package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SamplingAlgorithm = RandomSamplingAlgorithm{}

type RandomSamplingAlgorithm struct {
	Rule *RandomSamplingAlgorithmRule `json:"rule,omitempty"`
	Seed *int64                       `json:"seed,omitempty"`

	// Fields inherited from SamplingAlgorithm
}

var _ json.Marshaler = RandomSamplingAlgorithm{}

func (s RandomSamplingAlgorithm) MarshalJSON() ([]byte, error) {
	type wrapper RandomSamplingAlgorithm
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RandomSamplingAlgorithm: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RandomSamplingAlgorithm: %+v", err)
	}
	decoded["samplingAlgorithmType"] = "Random"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RandomSamplingAlgorithm: %+v", err)
	}

	return encoded, nil
}
