package storagesyncservicesresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageSyncServiceProperties struct {
	IncomingTrafficPolicy      *IncomingTrafficPolicy       `json:"incomingTrafficPolicy,omitempty"`
	LastOperationName          *string                      `json:"lastOperationName,omitempty"`
	LastWorkflowId             *string                      `json:"lastWorkflowId,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *string                      `json:"provisioningState,omitempty"`
	StorageSyncServiceStatus   *int64                       `json:"storageSyncServiceStatus,omitempty"`
	StorageSyncServiceUid      *string                      `json:"storageSyncServiceUid,omitempty"`
}
