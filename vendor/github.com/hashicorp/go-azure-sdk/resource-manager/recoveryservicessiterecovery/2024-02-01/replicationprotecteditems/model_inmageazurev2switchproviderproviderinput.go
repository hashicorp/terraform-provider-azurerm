package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SwitchProviderProviderSpecificInput = InMageAzureV2SwitchProviderProviderInput{}

type InMageAzureV2SwitchProviderProviderInput struct {
	TargetApplianceID string `json:"targetApplianceID"`
	TargetFabricID    string `json:"targetFabricID"`
	TargetVaultID     string `json:"targetVaultID"`

	// Fields inherited from SwitchProviderProviderSpecificInput
}

var _ json.Marshaler = InMageAzureV2SwitchProviderProviderInput{}

func (s InMageAzureV2SwitchProviderProviderInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageAzureV2SwitchProviderProviderInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageAzureV2SwitchProviderProviderInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageAzureV2SwitchProviderProviderInput: %+v", err)
	}
	decoded["instanceType"] = "InMageAzureV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageAzureV2SwitchProviderProviderInput: %+v", err)
	}

	return encoded, nil
}
