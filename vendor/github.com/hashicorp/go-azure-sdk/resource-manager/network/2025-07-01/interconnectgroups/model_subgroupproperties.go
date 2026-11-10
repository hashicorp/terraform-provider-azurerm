package interconnectgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubgroupProperties struct {
	InterconnectBlock  *SubResource       `json:"interconnectBlock,omitempty"`
	InternalSubgroupId *string            `json:"internalSubgroupId,omitempty"`
	ProvisioningState  *ProvisioningState `json:"provisioningState,omitempty"`
	VirtualMachines    *[]SubResource     `json:"virtualMachines,omitempty"`
}
