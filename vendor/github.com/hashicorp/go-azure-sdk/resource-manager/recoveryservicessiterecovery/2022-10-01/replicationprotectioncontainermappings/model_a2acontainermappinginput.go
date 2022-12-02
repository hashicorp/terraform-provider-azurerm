package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificContainerMappingInput = A2AContainerMappingInput{}

type A2AContainerMappingInput struct {
	AgentAutoUpdateStatus               *AgentAutoUpdateStatus               `json:"agentAutoUpdateStatus,omitempty"`
	AutomationAccountArmId              *string                              `json:"automationAccountArmId,omitempty"`
	AutomationAccountAuthenticationType *AutomationAccountAuthenticationType `json:"automationAccountAuthenticationType,omitempty"`

	// Fields inherited from ReplicationProviderSpecificContainerMappingInput
}

var _ json.Marshaler = A2AContainerMappingInput{}

func (s A2AContainerMappingInput) MarshalJSON() ([]byte, error) {
	type wrapper A2AContainerMappingInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2AContainerMappingInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2AContainerMappingInput: %+v", err)
	}
	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2AContainerMappingInput: %+v", err)
	}

	return encoded, nil
}
