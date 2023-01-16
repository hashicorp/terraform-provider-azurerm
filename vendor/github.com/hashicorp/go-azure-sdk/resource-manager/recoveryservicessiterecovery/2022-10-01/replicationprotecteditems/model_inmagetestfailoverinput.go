package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TestFailoverProviderSpecificInput = InMageTestFailoverInput{}

type InMageTestFailoverInput struct {
	RecoveryPointId   *string            `json:"recoveryPointId,omitempty"`
	RecoveryPointType *RecoveryPointType `json:"recoveryPointType,omitempty"`

	// Fields inherited from TestFailoverProviderSpecificInput
}

var _ json.Marshaler = InMageTestFailoverInput{}

func (s InMageTestFailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageTestFailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageTestFailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageTestFailoverInput: %+v", err)
	}
	decoded["instanceType"] = "InMage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageTestFailoverInput: %+v", err)
	}

	return encoded, nil
}
