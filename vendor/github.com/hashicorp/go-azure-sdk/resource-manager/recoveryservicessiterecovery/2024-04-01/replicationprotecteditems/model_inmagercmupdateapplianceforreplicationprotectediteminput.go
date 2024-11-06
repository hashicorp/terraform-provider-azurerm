package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UpdateApplianceForReplicationProtectedItemProviderSpecificInput = InMageRcmUpdateApplianceForReplicationProtectedItemInput{}

type InMageRcmUpdateApplianceForReplicationProtectedItemInput struct {
	RunAsAccountId *string `json:"runAsAccountId,omitempty"`

	// Fields inherited from UpdateApplianceForReplicationProtectedItemProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s InMageRcmUpdateApplianceForReplicationProtectedItemInput) UpdateApplianceForReplicationProtectedItemProviderSpecificInput() BaseUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl {
	return BaseUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageRcmUpdateApplianceForReplicationProtectedItemInput{}

func (s InMageRcmUpdateApplianceForReplicationProtectedItemInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmUpdateApplianceForReplicationProtectedItemInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmUpdateApplianceForReplicationProtectedItemInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmUpdateApplianceForReplicationProtectedItemInput: %+v", err)
	}

	decoded["instanceType"] = "InMageRcm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmUpdateApplianceForReplicationProtectedItemInput: %+v", err)
	}

	return encoded, nil
}
