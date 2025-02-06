package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapabilitiesResult struct {
	Features *[]string                      `json:"features,omitempty"`
	Quota    *QuotaCapability               `json:"quota,omitempty"`
	Regions  *map[string]RegionsCapability  `json:"regions,omitempty"`
	Versions *map[string]VersionsCapability `json:"versions,omitempty"`
}
