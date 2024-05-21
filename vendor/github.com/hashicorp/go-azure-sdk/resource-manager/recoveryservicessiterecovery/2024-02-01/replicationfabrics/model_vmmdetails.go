package replicationfabrics

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificDetails = VMmDetails{}

type VMmDetails struct {

	// Fields inherited from FabricSpecificDetails
}

var _ json.Marshaler = VMmDetails{}

func (s VMmDetails) MarshalJSON() ([]byte, error) {
	type wrapper VMmDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMmDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMmDetails: %+v", err)
	}
	decoded["instanceType"] = "VMM"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMmDetails: %+v", err)
	}

	return encoded, nil
}
