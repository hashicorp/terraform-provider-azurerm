package dataset

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataSet = ADLSGen1FileDataSet{}

type ADLSGen1FileDataSet struct {
	Properties ADLSGen1FileProperties `json:"properties"`

	// Fields inherited from DataSet
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

var _ json.Marshaler = ADLSGen1FileDataSet{}

func (s ADLSGen1FileDataSet) MarshalJSON() ([]byte, error) {
	type wrapper ADLSGen1FileDataSet
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ADLSGen1FileDataSet: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ADLSGen1FileDataSet: %+v", err)
	}
	decoded["kind"] = "AdlsGen1File"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ADLSGen1FileDataSet: %+v", err)
	}

	return encoded, nil
}
