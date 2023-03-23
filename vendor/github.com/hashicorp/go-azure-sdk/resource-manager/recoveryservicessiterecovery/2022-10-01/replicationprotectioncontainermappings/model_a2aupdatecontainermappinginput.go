package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificUpdateContainerMappingInput = A2AUpdateContainerMappingInput{}

type A2AUpdateContainerMappingInput struct {
	AgentAutoUpdateStatus               *AgentAutoUpdateStatus               `json:"agentAutoUpdateStatus,omitempty"`
	AutomationAccountArmId              *string                              `json:"automationAccountArmId,omitempty"`
	AutomationAccountAuthenticationType *AutomationAccountAuthenticationType `json:"automationAccountAuthenticationType,omitempty"`

	// Fields inherited from ReplicationProviderSpecificUpdateContainerMappingInput
}

var _ json.Marshaler = A2AUpdateContainerMappingInput{}

func (s A2AUpdateContainerMappingInput) MarshalJSON() ([]byte, error) {
	type wrapper A2AUpdateContainerMappingInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2AUpdateContainerMappingInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2AUpdateContainerMappingInput: %+v", err)
	}
	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2AUpdateContainerMappingInput: %+v", err)
	}

	return encoded, nil
}
