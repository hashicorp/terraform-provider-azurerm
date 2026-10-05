package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OsProfileUpdateLinuxConfiguration struct {
	ProvisionVMAgent       *bool `json:"provisionVMAgent,omitempty"`
	ProvisionVMConfigAgent *bool `json:"provisionVMConfigAgent,omitempty"`
}
