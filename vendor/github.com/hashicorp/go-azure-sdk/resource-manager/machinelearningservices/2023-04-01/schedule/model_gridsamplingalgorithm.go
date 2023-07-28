package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SamplingAlgorithm = GridSamplingAlgorithm{}

type GridSamplingAlgorithm struct {

	// Fields inherited from SamplingAlgorithm
}

var _ json.Marshaler = GridSamplingAlgorithm{}

func (s GridSamplingAlgorithm) MarshalJSON() ([]byte, error) {
	type wrapper GridSamplingAlgorithm
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling GridSamplingAlgorithm: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling GridSamplingAlgorithm: %+v", err)
	}
	decoded["samplingAlgorithmType"] = "Grid"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling GridSamplingAlgorithm: %+v", err)
	}

	return encoded, nil
}
