package alertprocessingrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Action = RemoveAllActionGroups{}

type RemoveAllActionGroups struct {

	// Fields inherited from Action

	ActionType ActionType `json:"actionType"`
}

func (s RemoveAllActionGroups) Action() BaseActionImpl {
	return BaseActionImpl{
		ActionType: s.ActionType,
	}
}

var _ json.Marshaler = RemoveAllActionGroups{}

func (s RemoveAllActionGroups) MarshalJSON() ([]byte, error) {
	type wrapper RemoveAllActionGroups
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RemoveAllActionGroups: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RemoveAllActionGroups: %+v", err)
	}

	decoded["actionType"] = "RemoveAllActionGroups"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RemoveAllActionGroups: %+v", err)
	}

	return encoded, nil
}
