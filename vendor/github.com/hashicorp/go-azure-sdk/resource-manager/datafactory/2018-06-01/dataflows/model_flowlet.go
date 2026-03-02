package dataflows

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataFlow = Flowlet{}

type Flowlet struct {
	TypeProperties *FlowletTypeProperties `json:"typeProperties,omitempty"`

	// Fields inherited from DataFlow

	Annotations *[]interface{}  `json:"annotations,omitempty"`
	Description *string         `json:"description,omitempty"`
	Folder      *DataFlowFolder `json:"folder,omitempty"`
	Type        string          `json:"type"`
}

func (s Flowlet) DataFlow() BaseDataFlowImpl {
	return BaseDataFlowImpl{
		Annotations: s.Annotations,
		Description: s.Description,
		Folder:      s.Folder,
		Type:        s.Type,
	}
}

var _ json.Marshaler = Flowlet{}

func (s Flowlet) MarshalJSON() ([]byte, error) {
	type wrapper Flowlet
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Flowlet: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Flowlet: %+v", err)
	}

	decoded["type"] = "Flowlet"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Flowlet: %+v", err)
	}

	return encoded, nil
}
