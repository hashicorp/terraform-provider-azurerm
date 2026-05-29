package alertprocessingrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Action = AddActionGroups{}

type AddActionGroups struct {
	ActionGroupIds []string `json:"actionGroupIds"`

	// Fields inherited from Action

	ActionType ActionType `json:"actionType"`
}

func (s AddActionGroups) Action() BaseActionImpl {
	return BaseActionImpl{
		ActionType: s.ActionType,
	}
}

var _ json.Marshaler = AddActionGroups{}

func (s AddActionGroups) MarshalJSON() ([]byte, error) {
	type wrapper AddActionGroups
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AddActionGroups: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AddActionGroups: %+v", err)
	}

	decoded["actionType"] = "AddActionGroups"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AddActionGroups: %+v", err)
	}

	return encoded, nil
}
