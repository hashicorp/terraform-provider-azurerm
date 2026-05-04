package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutputDataSource = AzureDataExplorerOutputDataSource{}

type AzureDataExplorerOutputDataSource struct {
	Properties *AzureDataExplorerOutputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from OutputDataSource

	Type string `json:"type"`
}

func (s AzureDataExplorerOutputDataSource) OutputDataSource() BaseOutputDataSourceImpl {
	return BaseOutputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AzureDataExplorerOutputDataSource{}

func (s AzureDataExplorerOutputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper AzureDataExplorerOutputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureDataExplorerOutputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureDataExplorerOutputDataSource: %+v", err)
	}

	decoded["type"] = "Microsoft.Kusto/clusters/databases"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureDataExplorerOutputDataSource: %+v", err)
	}

	return encoded, nil
}
