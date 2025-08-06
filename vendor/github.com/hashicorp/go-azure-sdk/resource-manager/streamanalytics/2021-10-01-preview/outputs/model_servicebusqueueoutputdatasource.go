package outputs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutputDataSource = ServiceBusQueueOutputDataSource{}

type ServiceBusQueueOutputDataSource struct {
	Properties *ServiceBusQueueOutputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from OutputDataSource

	Type string `json:"type"`
}

func (s ServiceBusQueueOutputDataSource) OutputDataSource() BaseOutputDataSourceImpl {
	return BaseOutputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ServiceBusQueueOutputDataSource{}

func (s ServiceBusQueueOutputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper ServiceBusQueueOutputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServiceBusQueueOutputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServiceBusQueueOutputDataSource: %+v", err)
	}

	decoded["type"] = "Microsoft.ServiceBus/Queue"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServiceBusQueueOutputDataSource: %+v", err)
	}

	return encoded, nil
}
