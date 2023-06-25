package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutputDataSource = AzureSynapseOutputDataSource{}

type AzureSynapseOutputDataSource struct {
	Properties *AzureSynapseDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from OutputDataSource
}

var _ json.Marshaler = AzureSynapseOutputDataSource{}

func (s AzureSynapseOutputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper AzureSynapseOutputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureSynapseOutputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureSynapseOutputDataSource: %+v", err)
	}
	decoded["type"] = "Microsoft.Sql/Server/DataWarehouse"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureSynapseOutputDataSource: %+v", err)
	}

	return encoded, nil
}
