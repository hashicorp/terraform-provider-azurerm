package inputs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamInputDataSource interface {
}

func unmarshalStreamInputDataSourceImplementation(input []byte) (StreamInputDataSource, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling StreamInputDataSource into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Microsoft.Storage/Blob") {
		var out BlobStreamInputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BlobStreamInputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.ServiceBus/EventHub") {
		var out EventHubStreamInputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHubStreamInputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.EventHub/EventHub") {
		var out EventHubV2StreamInputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHubV2StreamInputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GatewayMessageBus") {
		var out GatewayMessageBusStreamInputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GatewayMessageBusStreamInputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.Devices/IotHubs") {
		var out IoTHubStreamInputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IoTHubStreamInputDataSource: %+v", err)
		}
		return out, nil
	}

	type RawStreamInputDataSourceImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawStreamInputDataSourceImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
