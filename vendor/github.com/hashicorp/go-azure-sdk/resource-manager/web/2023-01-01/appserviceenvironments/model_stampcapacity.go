package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StampCapacity struct {
	AvailableCapacity              *int64              `json:"availableCapacity,omitempty"`
	ComputeMode                    *ComputeModeOptions `json:"computeMode,omitempty"`
	ExcludeFromCapacityAllocation  *bool               `json:"excludeFromCapacityAllocation,omitempty"`
	IsApplicableForAllComputeModes *bool               `json:"isApplicableForAllComputeModes,omitempty"`
	IsLinux                        *bool               `json:"isLinux,omitempty"`
	Name                           *string             `json:"name,omitempty"`
	SiteMode                       *string             `json:"siteMode,omitempty"`
	TotalCapacity                  *int64              `json:"totalCapacity,omitempty"`
	Unit                           *string             `json:"unit,omitempty"`
	WorkerSize                     *WorkerSizeOptions  `json:"workerSize,omitempty"`
	WorkerSizeId                   *int64              `json:"workerSizeId,omitempty"`
}
