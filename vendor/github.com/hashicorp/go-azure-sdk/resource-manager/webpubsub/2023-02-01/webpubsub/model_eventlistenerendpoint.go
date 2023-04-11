package webpubsub

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventListenerEndpoint interface {
}

func unmarshalEventListenerEndpointImplementation(input []byte) (EventListenerEndpoint, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EventListenerEndpoint into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "EventHub") {
		var out EventHubEndpoint
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHubEndpoint: %+v", err)
		}
		return out, nil
	}

	type RawEventListenerEndpointImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawEventListenerEndpointImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
