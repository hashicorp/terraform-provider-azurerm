package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineAgentInstanceView struct {
	ExtensionHandlers *[]VirtualMachineExtensionHandlerInstanceView `json:"extensionHandlers,omitempty"`
	Statuses          *[]InstanceViewStatus                         `json:"statuses,omitempty"`
	VMAgentVersion    *string                                       `json:"vmAgentVersion,omitempty"`
}
