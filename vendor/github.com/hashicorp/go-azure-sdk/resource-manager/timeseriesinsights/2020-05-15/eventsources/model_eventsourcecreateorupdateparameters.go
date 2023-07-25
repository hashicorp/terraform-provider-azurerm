package eventsources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSourceCreateOrUpdateParameters interface {
}

func unmarshalEventSourceCreateOrUpdateParametersImplementation(input []byte) (EventSourceCreateOrUpdateParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EventSourceCreateOrUpdateParameters into map[string]interface: %+v", err)
	}

	value, ok := temp["kind"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Microsoft.EventHub") {
		var out EventHubEventSourceCreateOrUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHubEventSourceCreateOrUpdateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.IoTHub") {
		var out IoTHubEventSourceCreateOrUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IoTHubEventSourceCreateOrUpdateParameters: %+v", err)
		}
		return out, nil
	}

	type RawEventSourceCreateOrUpdateParametersImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawEventSourceCreateOrUpdateParametersImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
