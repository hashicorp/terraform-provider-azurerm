package streamingjobs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamInputDataSource interface {
	StreamInputDataSource() BaseStreamInputDataSourceImpl
}

var _ StreamInputDataSource = BaseStreamInputDataSourceImpl{}

type BaseStreamInputDataSourceImpl struct {
	Type string `json:"type"`
}

func (s BaseStreamInputDataSourceImpl) StreamInputDataSource() BaseStreamInputDataSourceImpl {
	return s
}

var _ StreamInputDataSource = RawStreamInputDataSourceImpl{}

// RawStreamInputDataSourceImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawStreamInputDataSourceImpl struct {
	streamInputDataSource BaseStreamInputDataSourceImpl
	Type                  string
	Values                map[string]interface{}
}

func (s RawStreamInputDataSourceImpl) StreamInputDataSource() BaseStreamInputDataSourceImpl {
	return s.streamInputDataSource
}

func UnmarshalStreamInputDataSourceImplementation(input []byte) (StreamInputDataSource, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling StreamInputDataSource into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Microsoft.Storage/Blob") {
		var out BlobStreamInputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BlobStreamInputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.EventGrid/EventSubscriptions") {
		var out EventGridStreamInputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventGridStreamInputDataSource: %+v", err)
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

	if strings.EqualFold(value, "Raw") {
		var out RawStreamInputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RawStreamInputDataSource: %+v", err)
		}
		return out, nil
	}

	var parent BaseStreamInputDataSourceImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseStreamInputDataSourceImpl: %+v", err)
	}

	return RawStreamInputDataSourceImpl{
		streamInputDataSource: parent,
		Type:                  value,
		Values:                temp,
	}, nil

}
