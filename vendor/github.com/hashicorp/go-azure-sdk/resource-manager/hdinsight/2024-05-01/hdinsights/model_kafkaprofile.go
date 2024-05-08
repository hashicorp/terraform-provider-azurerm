package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KafkaProfile struct {
	ConnectivityEndpoints *KafkaConnectivityEndpoints `json:"connectivityEndpoints,omitempty"`
	DiskStorage           DiskStorageProfile          `json:"diskStorage"`
	EnableKRaft           *bool                       `json:"enableKRaft,omitempty"`
	EnablePublicEndpoints *bool                       `json:"enablePublicEndpoints,omitempty"`
	RemoteStorageUri      *string                     `json:"remoteStorageUri,omitempty"`
}
