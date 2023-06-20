package eventsources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSourceUpdateParameters interface {
}

func unmarshalEventSourceUpdateParametersImplementation(input []byte) (EventSourceUpdateParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EventSourceUpdateParameters into map[string]interface: %+v", err)
	}

	value, ok := temp["kind"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Microsoft.EventHub") {
		var out EventHubEventSourceUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHubEventSourceUpdateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.IoTHub") {
		var out IoTHubEventSourceUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IoTHubEventSourceUpdateParameters: %+v", err)
		}
		return out, nil
	}

	type RawEventSourceUpdateParametersImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawEventSourceUpdateParametersImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
