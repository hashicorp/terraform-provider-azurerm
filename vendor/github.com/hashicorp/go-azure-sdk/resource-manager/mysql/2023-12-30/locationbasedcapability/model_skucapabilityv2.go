package locationbasedcapability

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuCapabilityV2 struct {
	Name                      *string   `json:"name,omitempty"`
	SupportedHAMode           *[]string `json:"supportedHAMode,omitempty"`
	SupportedIops             *int64    `json:"supportedIops,omitempty"`
	SupportedMemoryPerVCoreMB *int64    `json:"supportedMemoryPerVCoreMB,omitempty"`
	SupportedZones            *[]string `json:"supportedZones,omitempty"`
	VCores                    *int64    `json:"vCores,omitempty"`
}
