package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFunctionEventSubscriptionDestinationProperties struct {
	DeliveryAttributeMappings     *[]DeliveryAttributeMapping `json:"deliveryAttributeMappings,omitempty"`
	MaxEventsPerBatch             *int64                      `json:"maxEventsPerBatch,omitempty"`
	PreferredBatchSizeInKilobytes *int64                      `json:"preferredBatchSizeInKilobytes,omitempty"`
	ResourceId                    *string                     `json:"resourceId,omitempty"`
}

var _ json.Unmarshaler = &AzureFunctionEventSubscriptionDestinationProperties{}

func (s *AzureFunctionEventSubscriptionDestinationProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		MaxEventsPerBatch             *int64  `json:"maxEventsPerBatch,omitempty"`
		PreferredBatchSizeInKilobytes *int64  `json:"preferredBatchSizeInKilobytes,omitempty"`
		ResourceId                    *string `json:"resourceId,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.MaxEventsPerBatch = decoded.MaxEventsPerBatch
	s.PreferredBatchSizeInKilobytes = decoded.PreferredBatchSizeInKilobytes
	s.ResourceId = decoded.ResourceId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureFunctionEventSubscriptionDestinationProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["deliveryAttributeMappings"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling DeliveryAttributeMappings into list []json.RawMessage: %+v", err)
		}

		output := make([]DeliveryAttributeMapping, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalDeliveryAttributeMappingImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'DeliveryAttributeMappings' for 'AzureFunctionEventSubscriptionDestinationProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.DeliveryAttributeMappings = &output
	}

	return nil
}
