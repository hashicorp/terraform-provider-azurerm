package contentkeypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ContentKeyPolicyConfiguration = ContentKeyPolicyFairPlayConfiguration{}

type ContentKeyPolicyFairPlayConfiguration struct {
	Ask                        string                                              `json:"ask"`
	FairPlayPfx                string                                              `json:"fairPlayPfx"`
	FairPlayPfxPassword        string                                              `json:"fairPlayPfxPassword"`
	OfflineRentalConfiguration *ContentKeyPolicyFairPlayOfflineRentalConfiguration `json:"offlineRentalConfiguration,omitempty"`
	RentalAndLeaseKeyType      ContentKeyPolicyFairPlayRentalAndLeaseKeyType       `json:"rentalAndLeaseKeyType"`
	RentalDuration             int64                                               `json:"rentalDuration"`

	// Fields inherited from ContentKeyPolicyConfiguration
}

var _ json.Marshaler = ContentKeyPolicyFairPlayConfiguration{}

func (s ContentKeyPolicyFairPlayConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper ContentKeyPolicyFairPlayConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContentKeyPolicyFairPlayConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicyFairPlayConfiguration: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.ContentKeyPolicyFairPlayConfiguration"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContentKeyPolicyFairPlayConfiguration: %+v", err)
	}

	return encoded, nil
}
