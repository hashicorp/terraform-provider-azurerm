package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMSizeCompatibilityFilterV2 struct {
	ClusterFlavors            *[]string   `json:"clusterFlavors,omitempty"`
	ClusterVersions           *[]string   `json:"clusterVersions,omitempty"`
	ComputeIsolationSupported *string     `json:"computeIsolationSupported,omitempty"`
	EspApplied                *string     `json:"espApplied,omitempty"`
	FilterMode                *FilterMode `json:"filterMode,omitempty"`
	NodeTypes                 *[]string   `json:"nodeTypes,omitempty"`
	OsType                    *[]OSType   `json:"osType,omitempty"`
	Regions                   *[]string   `json:"regions,omitempty"`
	VMSizes                   *[]string   `json:"vmSizes,omitempty"`
}
