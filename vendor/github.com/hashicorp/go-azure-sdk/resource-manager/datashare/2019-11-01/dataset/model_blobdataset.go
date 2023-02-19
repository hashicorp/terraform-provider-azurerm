package dataset

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataSet = BlobDataSet{}

type BlobDataSet struct {
	Properties BlobProperties `json:"properties"`

	// Fields inherited from DataSet
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

var _ json.Marshaler = BlobDataSet{}

func (s BlobDataSet) MarshalJSON() ([]byte, error) {
	type wrapper BlobDataSet
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BlobDataSet: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BlobDataSet: %+v", err)
	}
	decoded["kind"] = "Blob"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BlobDataSet: %+v", err)
	}

	return encoded, nil
}
