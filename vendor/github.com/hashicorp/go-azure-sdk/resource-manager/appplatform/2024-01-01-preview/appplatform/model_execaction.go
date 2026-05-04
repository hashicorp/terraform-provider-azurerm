package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProbeAction = ExecAction{}

type ExecAction struct {
	Command *[]string `json:"command,omitempty"`

	// Fields inherited from ProbeAction

	Type ProbeActionType `json:"type"`
}

func (s ExecAction) ProbeAction() BaseProbeActionImpl {
	return BaseProbeActionImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ExecAction{}

func (s ExecAction) MarshalJSON() ([]byte, error) {
	type wrapper ExecAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ExecAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ExecAction: %+v", err)
	}

	decoded["type"] = "ExecAction"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ExecAction: %+v", err)
	}

	return encoded, nil
}
