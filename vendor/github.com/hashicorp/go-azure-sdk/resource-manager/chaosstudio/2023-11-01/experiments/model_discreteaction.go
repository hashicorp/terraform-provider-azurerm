package experiments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Action = DiscreteAction{}

type DiscreteAction struct {
	Parameters []KeyValuePair `json:"parameters"`
	SelectorId string         `json:"selectorId"`

	// Fields inherited from Action

	Name string `json:"name"`
	Type string `json:"type"`
}

func (s DiscreteAction) Action() BaseActionImpl {
	return BaseActionImpl{
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = DiscreteAction{}

func (s DiscreteAction) MarshalJSON() ([]byte, error) {
	type wrapper DiscreteAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DiscreteAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DiscreteAction: %+v", err)
	}

	decoded["type"] = "discrete"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DiscreteAction: %+v", err)
	}

	return encoded, nil
}
