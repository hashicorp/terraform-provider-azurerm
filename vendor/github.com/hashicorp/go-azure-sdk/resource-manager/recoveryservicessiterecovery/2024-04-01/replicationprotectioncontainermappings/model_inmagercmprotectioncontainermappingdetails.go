package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProtectionContainerMappingProviderSpecificDetails = InMageRcmProtectionContainerMappingDetails{}

type InMageRcmProtectionContainerMappingDetails struct {
	EnableAgentAutoUpgrade *string `json:"enableAgentAutoUpgrade,omitempty"`

	// Fields inherited from ProtectionContainerMappingProviderSpecificDetails

	InstanceType string `json:"instanceType"`
}

func (s InMageRcmProtectionContainerMappingDetails) ProtectionContainerMappingProviderSpecificDetails() BaseProtectionContainerMappingProviderSpecificDetailsImpl {
	return BaseProtectionContainerMappingProviderSpecificDetailsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageRcmProtectionContainerMappingDetails{}

func (s InMageRcmProtectionContainerMappingDetails) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmProtectionContainerMappingDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmProtectionContainerMappingDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmProtectionContainerMappingDetails: %+v", err)
	}

	decoded["instanceType"] = "InMageRcm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmProtectionContainerMappingDetails: %+v", err)
	}

	return encoded, nil
}
