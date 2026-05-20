package dataset

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataSet = ADLSGen1FolderDataSet{}

type ADLSGen1FolderDataSet struct {
	Properties ADLSGen1FolderProperties `json:"properties"`

	// Fields inherited from DataSet

	Id   *string     `json:"id,omitempty"`
	Kind DataSetKind `json:"kind"`
	Name *string     `json:"name,omitempty"`
	Type *string     `json:"type,omitempty"`
}

func (s ADLSGen1FolderDataSet) DataSet() BaseDataSetImpl {
	return BaseDataSetImpl{
		Id:   s.Id,
		Kind: s.Kind,
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = ADLSGen1FolderDataSet{}

func (s ADLSGen1FolderDataSet) MarshalJSON() ([]byte, error) {
	type wrapper ADLSGen1FolderDataSet
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ADLSGen1FolderDataSet: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ADLSGen1FolderDataSet: %+v", err)
	}

	decoded["kind"] = "AdlsGen1Folder"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ADLSGen1FolderDataSet: %+v", err)
	}

	return encoded, nil
}
