package virtualnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkPropertiesStatusOperationStatus struct {
	Error       *VirtualNetworkPropertiesStatusOperationStatusError `json:"error,omitempty"`
	OperationId *string                                             `json:"operationId,omitempty"`
	Status      *string                                             `json:"status,omitempty"`
}
