package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ApplyRecoveryPointProviderSpecificInput = A2AApplyRecoveryPointInput{}

type A2AApplyRecoveryPointInput struct {

	// Fields inherited from ApplyRecoveryPointProviderSpecificInput
}

var _ json.Marshaler = A2AApplyRecoveryPointInput{}

func (s A2AApplyRecoveryPointInput) MarshalJSON() ([]byte, error) {
	type wrapper A2AApplyRecoveryPointInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2AApplyRecoveryPointInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2AApplyRecoveryPointInput: %+v", err)
	}
	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2AApplyRecoveryPointInput: %+v", err)
	}

	return encoded, nil
}
