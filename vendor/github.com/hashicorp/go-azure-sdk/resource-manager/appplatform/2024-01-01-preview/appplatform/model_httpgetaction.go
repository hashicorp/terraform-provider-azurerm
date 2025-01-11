package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProbeAction = HTTPGetAction{}

type HTTPGetAction struct {
	Path   *string         `json:"path,omitempty"`
	Scheme *HTTPSchemeType `json:"scheme,omitempty"`

	// Fields inherited from ProbeAction

	Type ProbeActionType `json:"type"`
}

func (s HTTPGetAction) ProbeAction() BaseProbeActionImpl {
	return BaseProbeActionImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = HTTPGetAction{}

func (s HTTPGetAction) MarshalJSON() ([]byte, error) {
	type wrapper HTTPGetAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HTTPGetAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HTTPGetAction: %+v", err)
	}

	decoded["type"] = "HTTPGetAction"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HTTPGetAction: %+v", err)
	}

	return encoded, nil
}
