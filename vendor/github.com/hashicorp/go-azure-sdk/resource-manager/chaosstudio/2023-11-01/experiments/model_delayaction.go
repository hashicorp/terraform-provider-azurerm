package experiments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Action = DelayAction{}

type DelayAction struct {
	Duration string `json:"duration"`

	// Fields inherited from Action

	Name string `json:"name"`
	Type string `json:"type"`
}

func (s DelayAction) Action() BaseActionImpl {
	return BaseActionImpl{
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = DelayAction{}

func (s DelayAction) MarshalJSON() ([]byte, error) {
	type wrapper DelayAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DelayAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DelayAction: %+v", err)
	}

	decoded["type"] = "delay"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DelayAction: %+v", err)
	}

	return encoded, nil
}
