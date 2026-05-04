package dataset

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataSet = ADLSGen2FolderDataSet{}

type ADLSGen2FolderDataSet struct {
	Properties ADLSGen2FolderProperties `json:"properties"`

	// Fields inherited from DataSet

	Id   *string     `json:"id,omitempty"`
	Kind DataSetKind `json:"kind"`
	Name *string     `json:"name,omitempty"`
	Type *string     `json:"type,omitempty"`
}

func (s ADLSGen2FolderDataSet) DataSet() BaseDataSetImpl {
	return BaseDataSetImpl{
		Id:   s.Id,
		Kind: s.Kind,
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = ADLSGen2FolderDataSet{}

func (s ADLSGen2FolderDataSet) MarshalJSON() ([]byte, error) {
	type wrapper ADLSGen2FolderDataSet
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ADLSGen2FolderDataSet: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ADLSGen2FolderDataSet: %+v", err)
	}

	decoded["kind"] = "AdlsGen2Folder"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ADLSGen2FolderDataSet: %+v", err)
	}

	return encoded, nil
}
