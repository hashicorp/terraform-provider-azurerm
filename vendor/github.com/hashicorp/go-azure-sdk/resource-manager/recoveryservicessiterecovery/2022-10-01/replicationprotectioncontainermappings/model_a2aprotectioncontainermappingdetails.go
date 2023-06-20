package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProtectionContainerMappingProviderSpecificDetails = A2AProtectionContainerMappingDetails{}

type A2AProtectionContainerMappingDetails struct {
	AgentAutoUpdateStatus               *AgentAutoUpdateStatus               `json:"agentAutoUpdateStatus,omitempty"`
	AutomationAccountArmId              *string                              `json:"automationAccountArmId,omitempty"`
	AutomationAccountAuthenticationType *AutomationAccountAuthenticationType `json:"automationAccountAuthenticationType,omitempty"`
	JobScheduleName                     *string                              `json:"jobScheduleName,omitempty"`
	ScheduleName                        *string                              `json:"scheduleName,omitempty"`

	// Fields inherited from ProtectionContainerMappingProviderSpecificDetails
}

var _ json.Marshaler = A2AProtectionContainerMappingDetails{}

func (s A2AProtectionContainerMappingDetails) MarshalJSON() ([]byte, error) {
	type wrapper A2AProtectionContainerMappingDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2AProtectionContainerMappingDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2AProtectionContainerMappingDetails: %+v", err)
	}
	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2AProtectionContainerMappingDetails: %+v", err)
	}

	return encoded, nil
}
