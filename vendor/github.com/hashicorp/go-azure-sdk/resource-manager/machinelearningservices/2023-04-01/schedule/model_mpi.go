package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DistributionConfiguration = Mpi{}

type Mpi struct {
	ProcessCountPerInstance *int64 `json:"processCountPerInstance,omitempty"`

	// Fields inherited from DistributionConfiguration
}

var _ json.Marshaler = Mpi{}

func (s Mpi) MarshalJSON() ([]byte, error) {
	type wrapper Mpi
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Mpi: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Mpi: %+v", err)
	}
	decoded["distributionType"] = "Mpi"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Mpi: %+v", err)
	}

	return encoded, nil
}
