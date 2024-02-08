package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HybridConnectionEventSubscriptionDestinationProperties struct {
	DeliveryAttributeMappings *[]DeliveryAttributeMapping `json:"deliveryAttributeMappings,omitempty"`
	ResourceId                *string                     `json:"resourceId,omitempty"`
}

var _ json.Unmarshaler = &HybridConnectionEventSubscriptionDestinationProperties{}

func (s *HybridConnectionEventSubscriptionDestinationProperties) UnmarshalJSON(bytes []byte) error {
	type alias HybridConnectionEventSubscriptionDestinationProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into HybridConnectionEventSubscriptionDestinationProperties: %+v", err)
	}

	s.ResourceId = decoded.ResourceId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling HybridConnectionEventSubscriptionDestinationProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["deliveryAttributeMappings"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling DeliveryAttributeMappings into list []json.RawMessage: %+v", err)
		}

		output := make([]DeliveryAttributeMapping, 0)
		for i, val := range listTemp {
			impl, err := unmarshalDeliveryAttributeMappingImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'DeliveryAttributeMappings' for 'HybridConnectionEventSubscriptionDestinationProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.DeliveryAttributeMappings = &output
	}
	return nil
}
