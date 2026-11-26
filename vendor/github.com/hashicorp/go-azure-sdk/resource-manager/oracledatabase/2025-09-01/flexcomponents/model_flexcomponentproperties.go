package flexcomponents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FlexComponentProperties struct {
	AvailableCoreCount         *int64        `json:"availableCoreCount,omitempty"`
	AvailableDbStorageInGbs    *int64        `json:"availableDbStorageInGbs,omitempty"`
	AvailableLocalStorageInGbs *int64        `json:"availableLocalStorageInGbs,omitempty"`
	AvailableMemoryInGbs       *int64        `json:"availableMemoryInGbs,omitempty"`
	ComputeModel               *string       `json:"computeModel,omitempty"`
	DescriptionSummary         *string       `json:"descriptionSummary,omitempty"`
	HardwareType               *HardwareType `json:"hardwareType,omitempty"`
	MinimumCoreCount           *int64        `json:"minimumCoreCount,omitempty"`
	RuntimeMinimumCoreCount    *int64        `json:"runtimeMinimumCoreCount,omitempty"`
	Shape                      *string       `json:"shape,omitempty"`
}
