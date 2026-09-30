package nodepools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodePoolProperties struct {
	ImageCacheLowerThreshold *int64             `json:"imageCacheLowerThreshold,omitempty"`
	ImageCacheUpperThreshold *int64             `json:"imageCacheUpperThreshold,omitempty"`
	MaxNodeCount             int64              `json:"maxNodeCount"`
	MinNodeCount             *int64             `json:"minNodeCount,omitempty"`
	OsDiskSizeGb             *int64             `json:"osDiskSizeGb,omitempty"`
	ProvisioningState        *ProvisioningState `json:"provisioningState,omitempty"`
	ScaleSetPriority         *ScaleSetPriority  `json:"scaleSetPriority,omitempty"`
	SubnetId                 string             `json:"subnetId"`
	VMSize                   VMSize             `json:"vmSize"`
}
