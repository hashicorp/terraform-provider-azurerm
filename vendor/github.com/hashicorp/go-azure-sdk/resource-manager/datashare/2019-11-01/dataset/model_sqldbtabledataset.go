package dataset

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataSet = SqlDBTableDataSet{}

type SqlDBTableDataSet struct {
	Properties *SqlDBTableProperties `json:"properties,omitempty"`

	// Fields inherited from DataSet
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

var _ json.Marshaler = SqlDBTableDataSet{}

func (s SqlDBTableDataSet) MarshalJSON() ([]byte, error) {
	type wrapper SqlDBTableDataSet
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SqlDBTableDataSet: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SqlDBTableDataSet: %+v", err)
	}
	decoded["kind"] = "SqlDBTable"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SqlDBTableDataSet: %+v", err)
	}

	return encoded, nil
}
