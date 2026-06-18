package integrationruntime

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ IntegrationRuntime = SelfHostedIntegrationRuntime{}

type SelfHostedIntegrationRuntime struct {
	TypeProperties *SelfHostedIntegrationRuntimeTypeProperties `json:"typeProperties,omitempty"`

	// Fields inherited from IntegrationRuntime

	Description *string                `json:"description,omitempty"`
	Type        IntegrationRuntimeType `json:"type"`
}

func (s SelfHostedIntegrationRuntime) IntegrationRuntime() BaseIntegrationRuntimeImpl {
	return BaseIntegrationRuntimeImpl{
		Description: s.Description,
		Type:        s.Type,
	}
}

var _ json.Marshaler = SelfHostedIntegrationRuntime{}

func (s SelfHostedIntegrationRuntime) MarshalJSON() ([]byte, error) {
	type wrapper SelfHostedIntegrationRuntime
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SelfHostedIntegrationRuntime: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SelfHostedIntegrationRuntime: %+v", err)
	}

	decoded["type"] = "SelfHosted"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SelfHostedIntegrationRuntime: %+v", err)
	}

	return encoded, nil
}
