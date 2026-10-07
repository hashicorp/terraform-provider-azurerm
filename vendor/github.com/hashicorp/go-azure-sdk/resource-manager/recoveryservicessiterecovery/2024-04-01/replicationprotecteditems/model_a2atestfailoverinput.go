package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TestFailoverProviderSpecificInput = A2ATestFailoverInput{}

type A2ATestFailoverInput struct {
	CloudServiceCreationOption *string `json:"cloudServiceCreationOption,omitempty"`
	RecoveryPointId            *string `json:"recoveryPointId,omitempty"`

	// Fields inherited from TestFailoverProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s A2ATestFailoverInput) TestFailoverProviderSpecificInput() BaseTestFailoverProviderSpecificInputImpl {
	return BaseTestFailoverProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = A2ATestFailoverInput{}

func (s A2ATestFailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper A2ATestFailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2ATestFailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2ATestFailoverInput: %+v", err)
	}

	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2ATestFailoverInput: %+v", err)
	}

	return encoded, nil
}
