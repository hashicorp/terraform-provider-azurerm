package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReferenceInputDataSource = FileReferenceInputDataSource{}

type FileReferenceInputDataSource struct {
	Properties *FileReferenceInputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from ReferenceInputDataSource

	Type string `json:"type"`
}

func (s FileReferenceInputDataSource) ReferenceInputDataSource() BaseReferenceInputDataSourceImpl {
	return BaseReferenceInputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = FileReferenceInputDataSource{}

func (s FileReferenceInputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper FileReferenceInputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FileReferenceInputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FileReferenceInputDataSource: %+v", err)
	}

	decoded["type"] = "File"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FileReferenceInputDataSource: %+v", err)
	}

	return encoded, nil
}
