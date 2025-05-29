package dataflows

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataFlow = MappingDataFlow{}

type MappingDataFlow struct {
	TypeProperties *MappingDataFlowTypeProperties `json:"typeProperties,omitempty"`

	// Fields inherited from DataFlow

	Annotations *[]interface{}  `json:"annotations,omitempty"`
	Description *string         `json:"description,omitempty"`
	Folder      *DataFlowFolder `json:"folder,omitempty"`
	Type        string          `json:"type"`
}

func (s MappingDataFlow) DataFlow() BaseDataFlowImpl {
	return BaseDataFlowImpl{
		Annotations: s.Annotations,
		Description: s.Description,
		Folder:      s.Folder,
		Type:        s.Type,
	}
}

var _ json.Marshaler = MappingDataFlow{}

func (s MappingDataFlow) MarshalJSON() ([]byte, error) {
	type wrapper MappingDataFlow
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MappingDataFlow: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MappingDataFlow: %+v", err)
	}

	decoded["type"] = "MappingDataFlow"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MappingDataFlow: %+v", err)
	}

	return encoded, nil
}
