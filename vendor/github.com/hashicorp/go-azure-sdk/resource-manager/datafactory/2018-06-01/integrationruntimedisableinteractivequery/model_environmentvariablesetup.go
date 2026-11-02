package integrationruntimedisableinteractivequery

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CustomSetupBase = EnvironmentVariableSetup{}

type EnvironmentVariableSetup struct {
	TypeProperties EnvironmentVariableSetupTypeProperties `json:"typeProperties"`

	// Fields inherited from CustomSetupBase

	Type string `json:"type"`
}

func (s EnvironmentVariableSetup) CustomSetupBase() BaseCustomSetupBaseImpl {
	return BaseCustomSetupBaseImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = EnvironmentVariableSetup{}

func (s EnvironmentVariableSetup) MarshalJSON() ([]byte, error) {
	type wrapper EnvironmentVariableSetup
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EnvironmentVariableSetup: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EnvironmentVariableSetup: %+v", err)
	}

	decoded["type"] = "EnvironmentVariableSetup"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EnvironmentVariableSetup: %+v", err)
	}

	return encoded, nil
}
