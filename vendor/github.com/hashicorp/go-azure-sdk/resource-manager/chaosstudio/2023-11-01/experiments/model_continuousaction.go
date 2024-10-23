package experiments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Action = ContinuousAction{}

type ContinuousAction struct {
	Duration   string         `json:"duration"`
	Parameters []KeyValuePair `json:"parameters"`
	SelectorId string         `json:"selectorId"`

	// Fields inherited from Action

	Name string `json:"name"`
	Type string `json:"type"`
}

func (s ContinuousAction) Action() BaseActionImpl {
	return BaseActionImpl{
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = ContinuousAction{}

func (s ContinuousAction) MarshalJSON() ([]byte, error) {
	type wrapper ContinuousAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContinuousAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContinuousAction: %+v", err)
	}

	decoded["type"] = "continuous"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContinuousAction: %+v", err)
	}

	return encoded, nil
}
