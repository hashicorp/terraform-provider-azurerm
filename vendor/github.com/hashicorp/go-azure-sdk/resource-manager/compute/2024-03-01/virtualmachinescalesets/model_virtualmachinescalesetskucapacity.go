package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetSkuCapacity struct {
	DefaultCapacity *int64                              `json:"defaultCapacity,omitempty"`
	Maximum         *int64                              `json:"maximum,omitempty"`
	Minimum         *int64                              `json:"minimum,omitempty"`
	ScaleType       *VirtualMachineScaleSetSkuScaleType `json:"scaleType,omitempty"`
}
