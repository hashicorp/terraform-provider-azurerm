package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetExtensionProfile struct {
	Extensions           *[]VirtualMachineScaleSetExtension `json:"extensions,omitempty"`
	ExtensionsTimeBudget *string                            `json:"extensionsTimeBudget,omitempty"`
}
