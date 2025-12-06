package dataflowendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowEndpointMqtt struct {
	Authentication       DataflowEndpointMqttAuthentication `json:"authentication"`
	ClientIdPrefix       *string                            `json:"clientIdPrefix,omitempty"`
	CloudEventAttributes *CloudEventAttributeType           `json:"cloudEventAttributes,omitempty"`
	Host                 *string                            `json:"host,omitempty"`
	KeepAliveSeconds     *int64                             `json:"keepAliveSeconds,omitempty"`
	MaxInflightMessages  *int64                             `json:"maxInflightMessages,omitempty"`
	Protocol             *BrokerProtocolType                `json:"protocol,omitempty"`
	Qos                  *int64                             `json:"qos,omitempty"`
	Retain               *MqttRetainType                    `json:"retain,omitempty"`
	SessionExpirySeconds *int64                             `json:"sessionExpirySeconds,omitempty"`
	Tls                  *TlsProperties                     `json:"tls,omitempty"`
}
