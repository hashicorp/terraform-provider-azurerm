package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuProfile struct {
	AllocationStrategy *AllocationStrategy `json:"allocationStrategy,omitempty"`
	VMSizes            *[]SkuProfileVMSize `json:"vmSizes,omitempty"`
}
