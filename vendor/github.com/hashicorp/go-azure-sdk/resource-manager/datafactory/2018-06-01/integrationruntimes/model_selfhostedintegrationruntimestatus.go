package integrationruntimes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ IntegrationRuntimeStatus = SelfHostedIntegrationRuntimeStatus{}

type SelfHostedIntegrationRuntimeStatus struct {
	TypeProperties SelfHostedIntegrationRuntimeStatusTypeProperties `json:"typeProperties"`

	// Fields inherited from IntegrationRuntimeStatus

	DataFactoryName *string                  `json:"dataFactoryName,omitempty"`
	State           *IntegrationRuntimeState `json:"state,omitempty"`
	Type            IntegrationRuntimeType   `json:"type"`
}

func (s SelfHostedIntegrationRuntimeStatus) IntegrationRuntimeStatus() BaseIntegrationRuntimeStatusImpl {
	return BaseIntegrationRuntimeStatusImpl{
		DataFactoryName: s.DataFactoryName,
		State:           s.State,
		Type:            s.Type,
	}
}

var _ json.Marshaler = SelfHostedIntegrationRuntimeStatus{}

func (s SelfHostedIntegrationRuntimeStatus) MarshalJSON() ([]byte, error) {
	type wrapper SelfHostedIntegrationRuntimeStatus
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SelfHostedIntegrationRuntimeStatus: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SelfHostedIntegrationRuntimeStatus: %+v", err)
	}

	decoded["type"] = "SelfHosted"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SelfHostedIntegrationRuntimeStatus: %+v", err)
	}

	return encoded, nil
}
