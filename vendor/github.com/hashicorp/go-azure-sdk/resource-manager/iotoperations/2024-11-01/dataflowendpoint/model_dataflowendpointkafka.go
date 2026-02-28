package dataflowendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowEndpointKafka struct {
	Authentication       DataflowEndpointKafkaAuthentication     `json:"authentication"`
	Batching             *DataflowEndpointKafkaBatching          `json:"batching,omitempty"`
	CloudEventAttributes *CloudEventAttributeType                `json:"cloudEventAttributes,omitempty"`
	Compression          *DataflowEndpointKafkaCompression       `json:"compression,omitempty"`
	ConsumerGroupId      *string                                 `json:"consumerGroupId,omitempty"`
	CopyMqttProperties   *OperationalMode                        `json:"copyMqttProperties,omitempty"`
	Host                 string                                  `json:"host"`
	KafkaAcks            *DataflowEndpointKafkaAcks              `json:"kafkaAcks,omitempty"`
	PartitionStrategy    *DataflowEndpointKafkaPartitionStrategy `json:"partitionStrategy,omitempty"`
	Tls                  *TlsProperties                          `json:"tls,omitempty"`
}
