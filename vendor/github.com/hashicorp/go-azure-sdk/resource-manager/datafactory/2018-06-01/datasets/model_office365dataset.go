package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Dataset = Office365Dataset{}

type Office365Dataset struct {
	TypeProperties Office365DatasetTypeProperties `json:"typeProperties"`

	// Fields inherited from Dataset

	Annotations       *[]interface{}                     `json:"annotations,omitempty"`
	Description       *string                            `json:"description,omitempty"`
	Folder            *DatasetFolder                     `json:"folder,omitempty"`
	LinkedServiceName LinkedServiceReference             `json:"linkedServiceName"`
	Parameters        *map[string]ParameterSpecification `json:"parameters,omitempty"`
	Schema            *interface{}                       `json:"schema,omitempty"`
	Structure         *interface{}                       `json:"structure,omitempty"`
	Type              string                             `json:"type"`
}

func (s Office365Dataset) Dataset() BaseDatasetImpl {
	return BaseDatasetImpl{
		Annotations:       s.Annotations,
		Description:       s.Description,
		Folder:            s.Folder,
		LinkedServiceName: s.LinkedServiceName,
		Parameters:        s.Parameters,
		Schema:            s.Schema,
		Structure:         s.Structure,
		Type:              s.Type,
	}
}

var _ json.Marshaler = Office365Dataset{}

func (s Office365Dataset) MarshalJSON() ([]byte, error) {
	type wrapper Office365Dataset
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Office365Dataset: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Office365Dataset: %+v", err)
	}

	decoded["type"] = "Office365Table"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Office365Dataset: %+v", err)
	}

	return encoded, nil
}
