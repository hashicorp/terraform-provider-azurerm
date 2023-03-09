package outputs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutputDataSource = AzureSqlDatabaseOutputDataSource{}

type AzureSqlDatabaseOutputDataSource struct {
	Properties *AzureSqlDatabaseDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from OutputDataSource
}

var _ json.Marshaler = AzureSqlDatabaseOutputDataSource{}

func (s AzureSqlDatabaseOutputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper AzureSqlDatabaseOutputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureSqlDatabaseOutputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureSqlDatabaseOutputDataSource: %+v", err)
	}
	decoded["type"] = "Microsoft.Sql/Server/Database"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureSqlDatabaseOutputDataSource: %+v", err)
	}

	return encoded, nil
}
