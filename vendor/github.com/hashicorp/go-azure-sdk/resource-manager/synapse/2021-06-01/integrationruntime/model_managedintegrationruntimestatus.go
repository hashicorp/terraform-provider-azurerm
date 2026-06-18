package integrationruntime

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ IntegrationRuntimeStatus = ManagedIntegrationRuntimeStatus{}

type ManagedIntegrationRuntimeStatus struct {
	TypeProperties ManagedIntegrationRuntimeStatusTypeProperties `json:"typeProperties"`

	// Fields inherited from IntegrationRuntimeStatus

	DataFactoryName *string                  `json:"dataFactoryName,omitempty"`
	State           *IntegrationRuntimeState `json:"state,omitempty"`
	Type            IntegrationRuntimeType   `json:"type"`
}

func (s ManagedIntegrationRuntimeStatus) IntegrationRuntimeStatus() BaseIntegrationRuntimeStatusImpl {
	return BaseIntegrationRuntimeStatusImpl{
		DataFactoryName: s.DataFactoryName,
		State:           s.State,
		Type:            s.Type,
	}
}

var _ json.Marshaler = ManagedIntegrationRuntimeStatus{}

func (s ManagedIntegrationRuntimeStatus) MarshalJSON() ([]byte, error) {
	type wrapper ManagedIntegrationRuntimeStatus
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ManagedIntegrationRuntimeStatus: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ManagedIntegrationRuntimeStatus: %+v", err)
	}

	decoded["type"] = "Managed"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ManagedIntegrationRuntimeStatus: %+v", err)
	}

	return encoded, nil
}
