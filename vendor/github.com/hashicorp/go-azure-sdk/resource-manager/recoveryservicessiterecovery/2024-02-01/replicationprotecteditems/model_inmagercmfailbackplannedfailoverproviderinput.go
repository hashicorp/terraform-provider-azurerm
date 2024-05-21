package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PlannedFailoverProviderSpecificFailoverInput = InMageRcmFailbackPlannedFailoverProviderInput{}

type InMageRcmFailbackPlannedFailoverProviderInput struct {
	RecoveryPointType InMageRcmFailbackRecoveryPointType `json:"recoveryPointType"`

	// Fields inherited from PlannedFailoverProviderSpecificFailoverInput
}

var _ json.Marshaler = InMageRcmFailbackPlannedFailoverProviderInput{}

func (s InMageRcmFailbackPlannedFailoverProviderInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmFailbackPlannedFailoverProviderInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmFailbackPlannedFailoverProviderInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmFailbackPlannedFailoverProviderInput: %+v", err)
	}
	decoded["instanceType"] = "InMageRcmFailback"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmFailbackPlannedFailoverProviderInput: %+v", err)
	}

	return encoded, nil
}
