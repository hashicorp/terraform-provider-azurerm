package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StreamInputDataSource = RawStreamInputDataSource{}

type RawStreamInputDataSource struct {
	Properties *RawInputDatasourceProperties `json:"properties,omitempty"`

	// Fields inherited from StreamInputDataSource

	Type string `json:"type"`
}

func (s RawStreamInputDataSource) StreamInputDataSource() BaseStreamInputDataSourceImpl {
	return BaseStreamInputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = RawStreamInputDataSource{}

func (s RawStreamInputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper RawStreamInputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RawStreamInputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RawStreamInputDataSource: %+v", err)
	}

	decoded["type"] = "Raw"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RawStreamInputDataSource: %+v", err)
	}

	return encoded, nil
}
