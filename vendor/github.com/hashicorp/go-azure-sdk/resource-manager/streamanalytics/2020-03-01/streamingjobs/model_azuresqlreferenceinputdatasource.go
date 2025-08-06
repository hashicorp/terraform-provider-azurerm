package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReferenceInputDataSource = AzureSqlReferenceInputDataSource{}

type AzureSqlReferenceInputDataSource struct {
	Properties *AzureSqlReferenceInputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from ReferenceInputDataSource

	Type string `json:"type"`
}

func (s AzureSqlReferenceInputDataSource) ReferenceInputDataSource() BaseReferenceInputDataSourceImpl {
	return BaseReferenceInputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AzureSqlReferenceInputDataSource{}

func (s AzureSqlReferenceInputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper AzureSqlReferenceInputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureSqlReferenceInputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureSqlReferenceInputDataSource: %+v", err)
	}

	decoded["type"] = "Microsoft.Sql/Server/Database"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureSqlReferenceInputDataSource: %+v", err)
	}

	return encoded, nil
}
