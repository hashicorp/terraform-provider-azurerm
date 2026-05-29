package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProbeAction = TCPSocketAction{}

type TCPSocketAction struct {

	// Fields inherited from ProbeAction

	Type ProbeActionType `json:"type"`
}

func (s TCPSocketAction) ProbeAction() BaseProbeActionImpl {
	return BaseProbeActionImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = TCPSocketAction{}

func (s TCPSocketAction) MarshalJSON() ([]byte, error) {
	type wrapper TCPSocketAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TCPSocketAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TCPSocketAction: %+v", err)
	}

	decoded["type"] = "TCPSocketAction"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TCPSocketAction: %+v", err)
	}

	return encoded, nil
}
