package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UnplannedFailoverProviderSpecificInput = InMageAzureV2UnplannedFailoverInput{}

type InMageAzureV2UnplannedFailoverInput struct {
	OsUpgradeVersion *string `json:"osUpgradeVersion,omitempty"`
	RecoveryPointId  *string `json:"recoveryPointId,omitempty"`

	// Fields inherited from UnplannedFailoverProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s InMageAzureV2UnplannedFailoverInput) UnplannedFailoverProviderSpecificInput() BaseUnplannedFailoverProviderSpecificInputImpl {
	return BaseUnplannedFailoverProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageAzureV2UnplannedFailoverInput{}

func (s InMageAzureV2UnplannedFailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageAzureV2UnplannedFailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageAzureV2UnplannedFailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageAzureV2UnplannedFailoverInput: %+v", err)
	}

	decoded["instanceType"] = "InMageAzureV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageAzureV2UnplannedFailoverInput: %+v", err)
	}

	return encoded, nil
}
