package dataset

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataSet = ADLSGen2FileDataSet{}

type ADLSGen2FileDataSet struct {
	Properties ADLSGen2FileProperties `json:"properties"`

	// Fields inherited from DataSet

	Id   *string     `json:"id,omitempty"`
	Kind DataSetKind `json:"kind"`
	Name *string     `json:"name,omitempty"`
	Type *string     `json:"type,omitempty"`
}

func (s ADLSGen2FileDataSet) DataSet() BaseDataSetImpl {
	return BaseDataSetImpl{
		Id:   s.Id,
		Kind: s.Kind,
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = ADLSGen2FileDataSet{}

func (s ADLSGen2FileDataSet) MarshalJSON() ([]byte, error) {
	type wrapper ADLSGen2FileDataSet
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ADLSGen2FileDataSet: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ADLSGen2FileDataSet: %+v", err)
	}

	decoded["kind"] = "AdlsGen2File"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ADLSGen2FileDataSet: %+v", err)
	}

	return encoded, nil
}
