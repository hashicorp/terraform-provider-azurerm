package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineExtensionHandlerInstanceView struct {
	Status             *InstanceViewStatus `json:"status,omitempty"`
	Type               *string             `json:"type,omitempty"`
	TypeHandlerVersion *string             `json:"typeHandlerVersion,omitempty"`
}
