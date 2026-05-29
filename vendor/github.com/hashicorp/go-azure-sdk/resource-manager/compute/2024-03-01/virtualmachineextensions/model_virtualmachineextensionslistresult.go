package virtualmachineextensions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineExtensionsListResult struct {
	Value *[]VirtualMachineExtension `json:"value,omitempty"`
}
