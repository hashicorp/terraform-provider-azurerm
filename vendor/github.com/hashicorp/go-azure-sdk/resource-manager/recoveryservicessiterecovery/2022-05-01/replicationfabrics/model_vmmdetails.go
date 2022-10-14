package replicationfabrics

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificDetails = VmmDetails{}

type VmmDetails struct {

	// Fields inherited from FabricSpecificDetails
}

var _ json.Marshaler = VmmDetails{}

func (s VmmDetails) MarshalJSON() ([]byte, error) {
	type wrapper VmmDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VmmDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VmmDetails: %+v", err)
	}
	decoded["instanceType"] = "VMM"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VmmDetails: %+v", err)
	}

	return encoded, nil
}
