package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutputDataSource = GatewayMessageBusOutputDataSource{}

type GatewayMessageBusOutputDataSource struct {
	Properties *GatewayMessageBusSourceProperties `json:"properties,omitempty"`

	// Fields inherited from OutputDataSource

	Type string `json:"type"`
}

func (s GatewayMessageBusOutputDataSource) OutputDataSource() BaseOutputDataSourceImpl {
	return BaseOutputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = GatewayMessageBusOutputDataSource{}

func (s GatewayMessageBusOutputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper GatewayMessageBusOutputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling GatewayMessageBusOutputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling GatewayMessageBusOutputDataSource: %+v", err)
	}

	decoded["type"] = "GatewayMessageBus"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling GatewayMessageBusOutputDataSource: %+v", err)
	}

	return encoded, nil
}
