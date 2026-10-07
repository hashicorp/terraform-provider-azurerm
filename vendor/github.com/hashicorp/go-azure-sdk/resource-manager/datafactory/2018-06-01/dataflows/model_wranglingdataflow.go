package dataflows

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataFlow = WranglingDataFlow{}

type WranglingDataFlow struct {
	TypeProperties *PowerQueryTypeProperties `json:"typeProperties,omitempty"`

	// Fields inherited from DataFlow

	Annotations *[]interface{}  `json:"annotations,omitempty"`
	Description *string         `json:"description,omitempty"`
	Folder      *DataFlowFolder `json:"folder,omitempty"`
	Type        string          `json:"type"`
}

func (s WranglingDataFlow) DataFlow() BaseDataFlowImpl {
	return BaseDataFlowImpl{
		Annotations: s.Annotations,
		Description: s.Description,
		Folder:      s.Folder,
		Type:        s.Type,
	}
}

var _ json.Marshaler = WranglingDataFlow{}

func (s WranglingDataFlow) MarshalJSON() ([]byte, error) {
	type wrapper WranglingDataFlow
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling WranglingDataFlow: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling WranglingDataFlow: %+v", err)
	}

	decoded["type"] = "WranglingDataFlow"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling WranglingDataFlow: %+v", err)
	}

	return encoded, nil
}
