package dataflowendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowEndpointProperties struct {
	DataExplorerSettings    *DataflowEndpointDataExplorer    `json:"dataExplorerSettings,omitempty"`
	DataLakeStorageSettings *DataflowEndpointDataLakeStorage `json:"dataLakeStorageSettings,omitempty"`
	EndpointType            EndpointType                     `json:"endpointType"`
	FabricOneLakeSettings   *DataflowEndpointFabricOneLake   `json:"fabricOneLakeSettings,omitempty"`
	KafkaSettings           *DataflowEndpointKafka           `json:"kafkaSettings,omitempty"`
	LocalStorageSettings    *DataflowEndpointLocalStorage    `json:"localStorageSettings,omitempty"`
	MqttSettings            *DataflowEndpointMqtt            `json:"mqttSettings,omitempty"`
	ProvisioningState       *ProvisioningState               `json:"provisioningState,omitempty"`
}
