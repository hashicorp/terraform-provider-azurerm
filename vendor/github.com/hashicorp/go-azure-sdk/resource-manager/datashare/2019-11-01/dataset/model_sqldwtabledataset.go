package dataset

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataSet = SqlDWTableDataSet{}

type SqlDWTableDataSet struct {
	Properties *SqlDWTableProperties `json:"properties,omitempty"`

	// Fields inherited from DataSet
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

var _ json.Marshaler = SqlDWTableDataSet{}

func (s SqlDWTableDataSet) MarshalJSON() ([]byte, error) {
	type wrapper SqlDWTableDataSet
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SqlDWTableDataSet: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SqlDWTableDataSet: %+v", err)
	}
	decoded["kind"] = "SqlDWTable"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SqlDWTableDataSet: %+v", err)
	}

	return encoded, nil
}
