package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UnplannedFailoverProviderSpecificInput = InMageUnplannedFailoverInput{}

type InMageUnplannedFailoverInput struct {
	RecoveryPointId   *string            `json:"recoveryPointId,omitempty"`
	RecoveryPointType *RecoveryPointType `json:"recoveryPointType,omitempty"`

	// Fields inherited from UnplannedFailoverProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s InMageUnplannedFailoverInput) UnplannedFailoverProviderSpecificInput() BaseUnplannedFailoverProviderSpecificInputImpl {
	return BaseUnplannedFailoverProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageUnplannedFailoverInput{}

func (s InMageUnplannedFailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageUnplannedFailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageUnplannedFailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageUnplannedFailoverInput: %+v", err)
	}

	decoded["instanceType"] = "InMage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageUnplannedFailoverInput: %+v", err)
	}

	return encoded, nil
}
