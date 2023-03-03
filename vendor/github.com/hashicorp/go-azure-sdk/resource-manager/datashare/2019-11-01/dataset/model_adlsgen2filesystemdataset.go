package dataset

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataSet = ADLSGen2FileSystemDataSet{}

type ADLSGen2FileSystemDataSet struct {
	Properties ADLSGen2FileSystemProperties `json:"properties"`

	// Fields inherited from DataSet
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

var _ json.Marshaler = ADLSGen2FileSystemDataSet{}

func (s ADLSGen2FileSystemDataSet) MarshalJSON() ([]byte, error) {
	type wrapper ADLSGen2FileSystemDataSet
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ADLSGen2FileSystemDataSet: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ADLSGen2FileSystemDataSet: %+v", err)
	}
	decoded["kind"] = "AdlsGen2FileSystem"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ADLSGen2FileSystemDataSet: %+v", err)
	}

	return encoded, nil
}
