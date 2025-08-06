package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MachineExtensionInstanceView struct {
	Name               *string                             `json:"name,omitempty"`
	Status             *MachineExtensionInstanceViewStatus `json:"status,omitempty"`
	Type               *string                             `json:"type,omitempty"`
	TypeHandlerVersion *string                             `json:"typeHandlerVersion,omitempty"`
}
