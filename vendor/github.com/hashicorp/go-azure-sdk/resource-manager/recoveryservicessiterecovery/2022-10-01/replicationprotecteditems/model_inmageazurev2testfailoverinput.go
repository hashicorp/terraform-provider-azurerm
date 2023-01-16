package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TestFailoverProviderSpecificInput = InMageAzureV2TestFailoverInput{}

type InMageAzureV2TestFailoverInput struct {
	RecoveryPointId *string `json:"recoveryPointId,omitempty"`

	// Fields inherited from TestFailoverProviderSpecificInput
}

var _ json.Marshaler = InMageAzureV2TestFailoverInput{}

func (s InMageAzureV2TestFailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageAzureV2TestFailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageAzureV2TestFailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageAzureV2TestFailoverInput: %+v", err)
	}
	decoded["instanceType"] = "InMageAzureV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageAzureV2TestFailoverInput: %+v", err)
	}

	return encoded, nil
}
