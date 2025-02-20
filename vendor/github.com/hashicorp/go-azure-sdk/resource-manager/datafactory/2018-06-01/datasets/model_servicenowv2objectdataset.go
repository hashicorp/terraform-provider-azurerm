package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Dataset = ServiceNowV2ObjectDataset{}

type ServiceNowV2ObjectDataset struct {
	TypeProperties *GenericDatasetTypeProperties `json:"typeProperties,omitempty"`

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

func (s ServiceNowV2ObjectDataset) Dataset() BaseDatasetImpl {
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

var _ json.Marshaler = ServiceNowV2ObjectDataset{}

func (s ServiceNowV2ObjectDataset) MarshalJSON() ([]byte, error) {
	type wrapper ServiceNowV2ObjectDataset
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServiceNowV2ObjectDataset: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServiceNowV2ObjectDataset: %+v", err)
	}

	decoded["type"] = "ServiceNowV2Object"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServiceNowV2ObjectDataset: %+v", err)
	}

	return encoded, nil
}
