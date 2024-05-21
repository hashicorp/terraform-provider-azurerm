package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UnplannedFailoverProviderSpecificInput = A2AUnplannedFailoverInput{}

type A2AUnplannedFailoverInput struct {
	CloudServiceCreationOption *string `json:"cloudServiceCreationOption,omitempty"`
	RecoveryPointId            *string `json:"recoveryPointId,omitempty"`

	// Fields inherited from UnplannedFailoverProviderSpecificInput
}

var _ json.Marshaler = A2AUnplannedFailoverInput{}

func (s A2AUnplannedFailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper A2AUnplannedFailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2AUnplannedFailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2AUnplannedFailoverInput: %+v", err)
	}
	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2AUnplannedFailoverInput: %+v", err)
	}

	return encoded, nil
}
