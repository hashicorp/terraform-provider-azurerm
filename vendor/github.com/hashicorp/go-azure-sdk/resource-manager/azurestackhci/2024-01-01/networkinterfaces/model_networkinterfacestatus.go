package networkinterfaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkInterfaceStatus struct {
	ErrorCode          *string                                   `json:"errorCode,omitempty"`
	ErrorMessage       *string                                   `json:"errorMessage,omitempty"`
	ProvisioningStatus *NetworkInterfaceStatusProvisioningStatus `json:"provisioningStatus,omitempty"`
}
