package locationbasedcapabilities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapabilityProperties struct {
	SupportedFlexibleServerEditions *[]ServerEditionCapability `json:"supportedFlexibleServerEditions,omitempty"`
	SupportedGeoBackupRegions       *[]string                  `json:"supportedGeoBackupRegions,omitempty"`
	SupportedHAMode                 *[]string                  `json:"supportedHAMode,omitempty"`
	Zone                            *string                    `json:"zone,omitempty"`
}
