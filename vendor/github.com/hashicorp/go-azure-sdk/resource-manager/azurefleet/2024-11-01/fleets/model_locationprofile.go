package fleets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationProfile struct {
	Location                      string                     `json:"location"`
	VirtualMachineProfileOverride *BaseVirtualMachineProfile `json:"virtualMachineProfileOverride,omitempty"`
}
