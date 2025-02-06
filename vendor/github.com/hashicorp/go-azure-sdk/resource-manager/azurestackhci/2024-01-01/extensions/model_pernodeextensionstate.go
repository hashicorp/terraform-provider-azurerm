package extensions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PerNodeExtensionState struct {
	Extension          *string                `json:"extension,omitempty"`
	InstanceView       *ExtensionInstanceView `json:"instanceView,omitempty"`
	Name               *string                `json:"name,omitempty"`
	State              *NodeExtensionState    `json:"state,omitempty"`
	TypeHandlerVersion *string                `json:"typeHandlerVersion,omitempty"`
}
