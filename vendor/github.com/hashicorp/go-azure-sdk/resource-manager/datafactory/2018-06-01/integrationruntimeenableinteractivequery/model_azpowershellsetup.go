package integrationruntimeenableinteractivequery

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CustomSetupBase = AzPowerShellSetup{}

type AzPowerShellSetup struct {
	TypeProperties AzPowerShellSetupTypeProperties `json:"typeProperties"`

	// Fields inherited from CustomSetupBase

	Type string `json:"type"`
}

func (s AzPowerShellSetup) CustomSetupBase() BaseCustomSetupBaseImpl {
	return BaseCustomSetupBaseImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AzPowerShellSetup{}

func (s AzPowerShellSetup) MarshalJSON() ([]byte, error) {
	type wrapper AzPowerShellSetup
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzPowerShellSetup: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzPowerShellSetup: %+v", err)
	}

	decoded["type"] = "AzPowerShellSetup"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzPowerShellSetup: %+v", err)
	}

	return encoded, nil
}
