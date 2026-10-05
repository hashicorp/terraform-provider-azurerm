package integrationruntimes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CustomSetupBase = ComponentSetup{}

type ComponentSetup struct {
	TypeProperties LicensedComponentSetupTypeProperties `json:"typeProperties"`

	// Fields inherited from CustomSetupBase

	Type string `json:"type"`
}

func (s ComponentSetup) CustomSetupBase() BaseCustomSetupBaseImpl {
	return BaseCustomSetupBaseImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ComponentSetup{}

func (s ComponentSetup) MarshalJSON() ([]byte, error) {
	type wrapper ComponentSetup
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ComponentSetup: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ComponentSetup: %+v", err)
	}

	decoded["type"] = "ComponentSetup"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ComponentSetup: %+v", err)
	}

	return encoded, nil
}
