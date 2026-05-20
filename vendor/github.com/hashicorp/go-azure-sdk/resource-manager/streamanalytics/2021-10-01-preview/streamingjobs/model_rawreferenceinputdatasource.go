package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReferenceInputDataSource = RawReferenceInputDataSource{}

type RawReferenceInputDataSource struct {
	Properties *RawInputDatasourceProperties `json:"properties,omitempty"`

	// Fields inherited from ReferenceInputDataSource

	Type string `json:"type"`
}

func (s RawReferenceInputDataSource) ReferenceInputDataSource() BaseReferenceInputDataSourceImpl {
	return BaseReferenceInputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = RawReferenceInputDataSource{}

func (s RawReferenceInputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper RawReferenceInputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RawReferenceInputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RawReferenceInputDataSource: %+v", err)
	}

	decoded["type"] = "Raw"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RawReferenceInputDataSource: %+v", err)
	}

	return encoded, nil
}
