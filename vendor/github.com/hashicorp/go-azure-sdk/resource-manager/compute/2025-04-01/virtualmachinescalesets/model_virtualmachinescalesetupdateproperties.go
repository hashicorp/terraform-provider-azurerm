package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetUpdateProperties struct {
	AdditionalCapabilities                 *AdditionalCapabilities                `json:"additionalCapabilities,omitempty"`
	AutomaticRepairsPolicy                 *AutomaticRepairsPolicy                `json:"automaticRepairsPolicy,omitempty"`
	DoNotRunExtensionsOnOverprovisionedVMs *bool                                  `json:"doNotRunExtensionsOnOverprovisionedVMs,omitempty"`
	Overprovision                          *bool                                  `json:"overprovision,omitempty"`
	PriorityMixPolicy                      *PriorityMixPolicy                     `json:"priorityMixPolicy,omitempty"`
	ProximityPlacementGroup                *SubResource                           `json:"proximityPlacementGroup,omitempty"`
	ResiliencyPolicy                       *ResiliencyPolicy                      `json:"resiliencyPolicy,omitempty"`
	ScaleInPolicy                          *ScaleInPolicy                         `json:"scaleInPolicy,omitempty"`
	SinglePlacementGroup                   *bool                                  `json:"singlePlacementGroup,omitempty"`
	SkuProfile                             *SkuProfile                            `json:"skuProfile,omitempty"`
	SpotRestorePolicy                      *SpotRestorePolicy                     `json:"spotRestorePolicy,omitempty"`
	UpgradePolicy                          *UpgradePolicy                         `json:"upgradePolicy,omitempty"`
	VirtualMachineProfile                  *VirtualMachineScaleSetUpdateVMProfile `json:"virtualMachineProfile,omitempty"`
	ZonalPlatformFaultDomainAlignMode      *ZonalPlatformFaultDomainAlignMode     `json:"zonalPlatformFaultDomainAlignMode,omitempty"`
}
