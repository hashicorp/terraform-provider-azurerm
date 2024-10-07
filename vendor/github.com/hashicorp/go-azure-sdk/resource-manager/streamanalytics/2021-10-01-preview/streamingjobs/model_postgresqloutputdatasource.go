package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutputDataSource = PostgreSQLOutputDataSource{}

type PostgreSQLOutputDataSource struct {
	Properties *PostgreSQLDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from OutputDataSource

	Type string `json:"type"`
}

func (s PostgreSQLOutputDataSource) OutputDataSource() BaseOutputDataSourceImpl {
	return BaseOutputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = PostgreSQLOutputDataSource{}

func (s PostgreSQLOutputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper PostgreSQLOutputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PostgreSQLOutputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PostgreSQLOutputDataSource: %+v", err)
	}

	decoded["type"] = "Microsoft.DBForPostgreSQL/servers/databases"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PostgreSQLOutputDataSource: %+v", err)
	}

	return encoded, nil
}
