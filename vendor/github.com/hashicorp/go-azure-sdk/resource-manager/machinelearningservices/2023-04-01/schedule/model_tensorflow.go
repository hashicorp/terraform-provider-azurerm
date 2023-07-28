package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DistributionConfiguration = TensorFlow{}

type TensorFlow struct {
	ParameterServerCount *int64 `json:"parameterServerCount,omitempty"`
	WorkerCount          *int64 `json:"workerCount,omitempty"`

	// Fields inherited from DistributionConfiguration
}

var _ json.Marshaler = TensorFlow{}

func (s TensorFlow) MarshalJSON() ([]byte, error) {
	type wrapper TensorFlow
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TensorFlow: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TensorFlow: %+v", err)
	}
	decoded["distributionType"] = "TensorFlow"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TensorFlow: %+v", err)
	}

	return encoded, nil
}
