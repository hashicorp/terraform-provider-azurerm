package replicationprotectioncontainers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificContainerCreationInput = A2AContainerCreationInput{}

type A2AContainerCreationInput struct {

	// Fields inherited from ReplicationProviderSpecificContainerCreationInput
}

var _ json.Marshaler = A2AContainerCreationInput{}

func (s A2AContainerCreationInput) MarshalJSON() ([]byte, error) {
	type wrapper A2AContainerCreationInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2AContainerCreationInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2AContainerCreationInput: %+v", err)
	}
	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2AContainerCreationInput: %+v", err)
	}

	return encoded, nil
}
