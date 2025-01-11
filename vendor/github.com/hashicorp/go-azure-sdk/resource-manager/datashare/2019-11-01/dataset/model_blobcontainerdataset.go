package dataset

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataSet = BlobContainerDataSet{}

type BlobContainerDataSet struct {
	Properties BlobContainerProperties `json:"properties"`

	// Fields inherited from DataSet

	Id   *string     `json:"id,omitempty"`
	Kind DataSetKind `json:"kind"`
	Name *string     `json:"name,omitempty"`
	Type *string     `json:"type,omitempty"`
}

func (s BlobContainerDataSet) DataSet() BaseDataSetImpl {
	return BaseDataSetImpl{
		Id:   s.Id,
		Kind: s.Kind,
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = BlobContainerDataSet{}

func (s BlobContainerDataSet) MarshalJSON() ([]byte, error) {
	type wrapper BlobContainerDataSet
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BlobContainerDataSet: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BlobContainerDataSet: %+v", err)
	}

	decoded["kind"] = "Container"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BlobContainerDataSet: %+v", err)
	}

	return encoded, nil
}
