package topics

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicProperties struct {
	DataResidencyBoundary      *DataResidencyBoundary       `json:"dataResidencyBoundary,omitempty"`
	DisableLocalAuth           *bool                        `json:"disableLocalAuth,omitempty"`
	Endpoint                   *string                      `json:"endpoint,omitempty"`
	InboundIPRules             *[]InboundIPRule             `json:"inboundIpRules,omitempty"`
	InputSchema                *InputSchema                 `json:"inputSchema,omitempty"`
	InputSchemaMapping         InputSchemaMapping           `json:"inputSchemaMapping"`
	MetricResourceId           *string                      `json:"metricResourceId,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *TopicProvisioningState      `json:"provisioningState,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccess         `json:"publicNetworkAccess,omitempty"`
}

var _ json.Unmarshaler = &TopicProperties{}

func (s *TopicProperties) UnmarshalJSON(bytes []byte) error {
	type alias TopicProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into TopicProperties: %+v", err)
	}

	s.DataResidencyBoundary = decoded.DataResidencyBoundary
	s.DisableLocalAuth = decoded.DisableLocalAuth
	s.Endpoint = decoded.Endpoint
	s.InboundIPRules = decoded.InboundIPRules
	s.InputSchema = decoded.InputSchema
	s.MetricResourceId = decoded.MetricResourceId
	s.PrivateEndpointConnections = decoded.PrivateEndpointConnections
	s.ProvisioningState = decoded.ProvisioningState
	s.PublicNetworkAccess = decoded.PublicNetworkAccess

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling TopicProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["inputSchemaMapping"]; ok {
		impl, err := unmarshalInputSchemaMappingImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'InputSchemaMapping' for 'TopicProperties': %+v", err)
		}
		s.InputSchemaMapping = impl
	}
	return nil
}
