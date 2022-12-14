package outputs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutputDataSource = PowerBIOutputDataSource{}

type PowerBIOutputDataSource struct {
	Properties *PowerBIOutputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from OutputDataSource
}

var _ json.Marshaler = PowerBIOutputDataSource{}

func (s PowerBIOutputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper PowerBIOutputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PowerBIOutputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PowerBIOutputDataSource: %+v", err)
	}
	decoded["type"] = "PowerBI"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PowerBIOutputDataSource: %+v", err)
	}

	return encoded, nil
}
