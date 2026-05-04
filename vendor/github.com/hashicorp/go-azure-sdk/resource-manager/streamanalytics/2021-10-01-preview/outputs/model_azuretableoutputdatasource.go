package outputs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutputDataSource = AzureTableOutputDataSource{}

type AzureTableOutputDataSource struct {
	Properties *AzureTableOutputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from OutputDataSource

	Type string `json:"type"`
}

func (s AzureTableOutputDataSource) OutputDataSource() BaseOutputDataSourceImpl {
	return BaseOutputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AzureTableOutputDataSource{}

func (s AzureTableOutputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper AzureTableOutputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureTableOutputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureTableOutputDataSource: %+v", err)
	}

	decoded["type"] = "Microsoft.Storage/Table"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureTableOutputDataSource: %+v", err)
	}

	return encoded, nil
}
