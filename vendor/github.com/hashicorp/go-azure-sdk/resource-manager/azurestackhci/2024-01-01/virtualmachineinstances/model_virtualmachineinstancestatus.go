package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineInstanceStatus struct {
	ErrorCode          *string                                         `json:"errorCode,omitempty"`
	ErrorMessage       *string                                         `json:"errorMessage,omitempty"`
	PowerState         *PowerStateEnum                                 `json:"powerState,omitempty"`
	ProvisioningStatus *VirtualMachineInstanceStatusProvisioningStatus `json:"provisioningStatus,omitempty"`
}
