package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StreamInputDataSource = IoTHubStreamInputDataSource{}

type IoTHubStreamInputDataSource struct {
	Properties *IoTHubStreamInputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from StreamInputDataSource

	Type string `json:"type"`
}

func (s IoTHubStreamInputDataSource) StreamInputDataSource() BaseStreamInputDataSourceImpl {
	return BaseStreamInputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = IoTHubStreamInputDataSource{}

func (s IoTHubStreamInputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper IoTHubStreamInputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling IoTHubStreamInputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling IoTHubStreamInputDataSource: %+v", err)
	}

	decoded["type"] = "Microsoft.Devices/IotHubs"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling IoTHubStreamInputDataSource: %+v", err)
	}

	return encoded, nil
}
