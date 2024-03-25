package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineInstanceStatusProvisioningStatus struct {
	OperationId *string `json:"operationId,omitempty"`
	Status      *Status `json:"status,omitempty"`
}
