package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutputDataSource = DocumentDbOutputDataSource{}

type DocumentDbOutputDataSource struct {
	Properties *DocumentDbOutputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from OutputDataSource

	Type string `json:"type"`
}

func (s DocumentDbOutputDataSource) OutputDataSource() BaseOutputDataSourceImpl {
	return BaseOutputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = DocumentDbOutputDataSource{}

func (s DocumentDbOutputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper DocumentDbOutputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DocumentDbOutputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DocumentDbOutputDataSource: %+v", err)
	}

	decoded["type"] = "Microsoft.Storage/DocumentDB"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DocumentDbOutputDataSource: %+v", err)
	}

	return encoded, nil
}
