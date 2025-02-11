package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutputDataSource = ServiceBusTopicOutputDataSource{}

type ServiceBusTopicOutputDataSource struct {
	Properties *ServiceBusTopicOutputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from OutputDataSource

	Type string `json:"type"`
}

func (s ServiceBusTopicOutputDataSource) OutputDataSource() BaseOutputDataSourceImpl {
	return BaseOutputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ServiceBusTopicOutputDataSource{}

func (s ServiceBusTopicOutputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper ServiceBusTopicOutputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServiceBusTopicOutputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServiceBusTopicOutputDataSource: %+v", err)
	}

	decoded["type"] = "Microsoft.ServiceBus/Topic"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServiceBusTopicOutputDataSource: %+v", err)
	}

	return encoded, nil
}
