package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificUpdateContainerMappingInput = InMageRcmUpdateContainerMappingInput{}

type InMageRcmUpdateContainerMappingInput struct {
	EnableAgentAutoUpgrade string `json:"enableAgentAutoUpgrade"`

	// Fields inherited from ReplicationProviderSpecificUpdateContainerMappingInput

	InstanceType string `json:"instanceType"`
}

func (s InMageRcmUpdateContainerMappingInput) ReplicationProviderSpecificUpdateContainerMappingInput() BaseReplicationProviderSpecificUpdateContainerMappingInputImpl {
	return BaseReplicationProviderSpecificUpdateContainerMappingInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageRcmUpdateContainerMappingInput{}

func (s InMageRcmUpdateContainerMappingInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmUpdateContainerMappingInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmUpdateContainerMappingInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmUpdateContainerMappingInput: %+v", err)
	}

	decoded["instanceType"] = "InMageRcm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmUpdateContainerMappingInput: %+v", err)
	}

	return encoded, nil
}
