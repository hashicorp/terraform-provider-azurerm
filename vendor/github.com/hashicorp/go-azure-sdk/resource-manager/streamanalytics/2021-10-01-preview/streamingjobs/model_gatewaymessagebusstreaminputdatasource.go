package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StreamInputDataSource = GatewayMessageBusStreamInputDataSource{}

type GatewayMessageBusStreamInputDataSource struct {
	Properties *GatewayMessageBusSourceProperties `json:"properties,omitempty"`

	// Fields inherited from StreamInputDataSource

	Type string `json:"type"`
}

func (s GatewayMessageBusStreamInputDataSource) StreamInputDataSource() BaseStreamInputDataSourceImpl {
	return BaseStreamInputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = GatewayMessageBusStreamInputDataSource{}

func (s GatewayMessageBusStreamInputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper GatewayMessageBusStreamInputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling GatewayMessageBusStreamInputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling GatewayMessageBusStreamInputDataSource: %+v", err)
	}

	decoded["type"] = "GatewayMessageBus"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling GatewayMessageBusStreamInputDataSource: %+v", err)
	}

	return encoded, nil
}
