package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DistributionConfiguration = PyTorch{}

type PyTorch struct {
	ProcessCountPerInstance *int64 `json:"processCountPerInstance,omitempty"`

	// Fields inherited from DistributionConfiguration
}

var _ json.Marshaler = PyTorch{}

func (s PyTorch) MarshalJSON() ([]byte, error) {
	type wrapper PyTorch
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PyTorch: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PyTorch: %+v", err)
	}
	decoded["distributionType"] = "PyTorch"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PyTorch: %+v", err)
	}

	return encoded, nil
}
