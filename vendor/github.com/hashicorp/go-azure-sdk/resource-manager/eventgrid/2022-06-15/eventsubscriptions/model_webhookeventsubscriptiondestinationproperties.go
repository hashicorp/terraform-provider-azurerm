package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebHookEventSubscriptionDestinationProperties struct {
	AzureActiveDirectoryApplicationIdOrUri *string                     `json:"azureActiveDirectoryApplicationIdOrUri,omitempty"`
	AzureActiveDirectoryTenantId           *string                     `json:"azureActiveDirectoryTenantId,omitempty"`
	DeliveryAttributeMappings              *[]DeliveryAttributeMapping `json:"deliveryAttributeMappings,omitempty"`
	EndpointBaseURL                        *string                     `json:"endpointBaseUrl,omitempty"`
	EndpointURL                            *string                     `json:"endpointUrl,omitempty"`
	MaxEventsPerBatch                      *int64                      `json:"maxEventsPerBatch,omitempty"`
	PreferredBatchSizeInKilobytes          *int64                      `json:"preferredBatchSizeInKilobytes,omitempty"`
}

var _ json.Unmarshaler = &WebHookEventSubscriptionDestinationProperties{}

func (s *WebHookEventSubscriptionDestinationProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AzureActiveDirectoryApplicationIdOrUri *string `json:"azureActiveDirectoryApplicationIdOrUri,omitempty"`
		AzureActiveDirectoryTenantId           *string `json:"azureActiveDirectoryTenantId,omitempty"`
		EndpointBaseURL                        *string `json:"endpointBaseUrl,omitempty"`
		EndpointURL                            *string `json:"endpointUrl,omitempty"`
		MaxEventsPerBatch                      *int64  `json:"maxEventsPerBatch,omitempty"`
		PreferredBatchSizeInKilobytes          *int64  `json:"preferredBatchSizeInKilobytes,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AzureActiveDirectoryApplicationIdOrUri = decoded.AzureActiveDirectoryApplicationIdOrUri
	s.AzureActiveDirectoryTenantId = decoded.AzureActiveDirectoryTenantId
	s.EndpointBaseURL = decoded.EndpointBaseURL
	s.EndpointURL = decoded.EndpointURL
	s.MaxEventsPerBatch = decoded.MaxEventsPerBatch
	s.PreferredBatchSizeInKilobytes = decoded.PreferredBatchSizeInKilobytes

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling WebHookEventSubscriptionDestinationProperties into map[string]json.RawMessage: %+v", err)
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
				return fmt.Errorf("unmarshaling index %d field 'DeliveryAttributeMappings' for 'WebHookEventSubscriptionDestinationProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.DeliveryAttributeMappings = &output
	}

	return nil
}
