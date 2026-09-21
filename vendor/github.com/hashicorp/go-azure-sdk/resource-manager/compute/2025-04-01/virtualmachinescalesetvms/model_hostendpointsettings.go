package virtualmachinescalesetvms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HostEndpointSettings struct {
	InVMAccessControlProfileReferenceId *string `json:"inVMAccessControlProfileReferenceId,omitempty"`
	Mode                                *Modes  `json:"mode,omitempty"`
}
