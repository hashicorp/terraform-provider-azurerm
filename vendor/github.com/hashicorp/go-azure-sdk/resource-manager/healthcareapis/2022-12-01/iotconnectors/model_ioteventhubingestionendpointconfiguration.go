package iotconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IotEventHubIngestionEndpointConfiguration struct {
	ConsumerGroup                   *string `json:"consumerGroup,omitempty"`
	EventHubName                    *string `json:"eventHubName,omitempty"`
	FullyQualifiedEventHubNamespace *string `json:"fullyQualifiedEventHubNamespace,omitempty"`
}
