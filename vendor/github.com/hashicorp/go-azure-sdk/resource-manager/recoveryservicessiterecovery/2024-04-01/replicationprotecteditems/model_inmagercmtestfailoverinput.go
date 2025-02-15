package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TestFailoverProviderSpecificInput = InMageRcmTestFailoverInput{}

type InMageRcmTestFailoverInput struct {
	NetworkId       *string `json:"networkId,omitempty"`
	RecoveryPointId *string `json:"recoveryPointId,omitempty"`

	// Fields inherited from TestFailoverProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s InMageRcmTestFailoverInput) TestFailoverProviderSpecificInput() BaseTestFailoverProviderSpecificInputImpl {
	return BaseTestFailoverProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageRcmTestFailoverInput{}

func (s InMageRcmTestFailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmTestFailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmTestFailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmTestFailoverInput: %+v", err)
	}

	decoded["instanceType"] = "InMageRcm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmTestFailoverInput: %+v", err)
	}

	return encoded, nil
}
