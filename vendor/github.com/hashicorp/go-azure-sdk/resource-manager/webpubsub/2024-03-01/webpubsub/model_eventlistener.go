package webpubsub

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventListener struct {
	Endpoint EventListenerEndpoint `json:"endpoint"`
	Filter   EventListenerFilter   `json:"filter"`
}

var _ json.Unmarshaler = &EventListener{}

func (s *EventListener) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling EventListener into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["endpoint"]; ok {
		impl, err := UnmarshalEventListenerEndpointImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Endpoint' for 'EventListener': %+v", err)
		}
		s.Endpoint = impl
	}

	if v, ok := temp["filter"]; ok {
		impl, err := UnmarshalEventListenerFilterImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Filter' for 'EventListener': %+v", err)
		}
		s.Filter = impl
	}

	return nil
}
