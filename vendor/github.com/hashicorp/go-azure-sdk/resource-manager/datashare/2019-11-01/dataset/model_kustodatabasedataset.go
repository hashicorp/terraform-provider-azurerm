package dataset

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataSet = KustoDatabaseDataSet{}

type KustoDatabaseDataSet struct {
	Properties KustoDatabaseDataSetProperties `json:"properties"`

	// Fields inherited from DataSet

	Id   *string     `json:"id,omitempty"`
	Kind DataSetKind `json:"kind"`
	Name *string     `json:"name,omitempty"`
	Type *string     `json:"type,omitempty"`
}

func (s KustoDatabaseDataSet) DataSet() BaseDataSetImpl {
	return BaseDataSetImpl{
		Id:   s.Id,
		Kind: s.Kind,
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = KustoDatabaseDataSet{}

func (s KustoDatabaseDataSet) MarshalJSON() ([]byte, error) {
	type wrapper KustoDatabaseDataSet
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KustoDatabaseDataSet: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KustoDatabaseDataSet: %+v", err)
	}

	decoded["kind"] = "KustoDatabase"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KustoDatabaseDataSet: %+v", err)
	}

	return encoded, nil
}
