package virtualmachineimages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineImageProperties struct {
	Architecture                 *ArchitectureTypes            `json:"architecture,omitempty"`
	AutomaticOSUpgradeProperties *AutomaticOSUpgradeProperties `json:"automaticOSUpgradeProperties,omitempty"`
	DataDiskImages               *[]DataDiskImage              `json:"dataDiskImages,omitempty"`
	Disallowed                   *DisallowedConfiguration      `json:"disallowed,omitempty"`
	Features                     *[]VirtualMachineImageFeature `json:"features,omitempty"`
	HyperVGeneration             *HyperVGenerationTypes        `json:"hyperVGeneration,omitempty"`
	ImageDeprecationStatus       *ImageDeprecationStatus       `json:"imageDeprecationStatus,omitempty"`
	OsDiskImage                  *OSDiskImage                  `json:"osDiskImage,omitempty"`
	Plan                         *PurchasePlan                 `json:"plan,omitempty"`
}
