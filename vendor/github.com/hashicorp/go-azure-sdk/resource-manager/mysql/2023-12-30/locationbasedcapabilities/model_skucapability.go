package locationbasedcapabilities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuCapability struct {
	Name                      *string `json:"name,omitempty"`
	SupportedIops             *int64  `json:"supportedIops,omitempty"`
	SupportedMemoryPerVCoreMB *int64  `json:"supportedMemoryPerVCoreMB,omitempty"`
	VCores                    *int64  `json:"vCores,omitempty"`
}
