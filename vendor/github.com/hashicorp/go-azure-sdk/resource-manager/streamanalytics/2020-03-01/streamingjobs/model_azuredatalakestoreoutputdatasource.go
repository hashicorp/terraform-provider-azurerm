package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutputDataSource = AzureDataLakeStoreOutputDataSource{}

type AzureDataLakeStoreOutputDataSource struct {
	Properties *AzureDataLakeStoreOutputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from OutputDataSource

	Type string `json:"type"`
}

func (s AzureDataLakeStoreOutputDataSource) OutputDataSource() BaseOutputDataSourceImpl {
	return BaseOutputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AzureDataLakeStoreOutputDataSource{}

func (s AzureDataLakeStoreOutputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper AzureDataLakeStoreOutputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureDataLakeStoreOutputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureDataLakeStoreOutputDataSource: %+v", err)
	}

	decoded["type"] = "Microsoft.DataLake/Accounts"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureDataLakeStoreOutputDataSource: %+v", err)
	}

	return encoded, nil
}
