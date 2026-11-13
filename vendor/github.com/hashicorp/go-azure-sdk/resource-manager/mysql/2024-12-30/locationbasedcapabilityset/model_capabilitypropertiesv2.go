package locationbasedcapabilityset

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapabilityPropertiesV2 struct {
	SupportedFeatures               *[]FeatureProperty           `json:"supportedFeatures,omitempty"`
	SupportedFlexibleServerEditions *[]ServerEditionCapabilityV2 `json:"supportedFlexibleServerEditions,omitempty"`
	SupportedGeoBackupRegions       *[]string                    `json:"supportedGeoBackupRegions,omitempty"`
	SupportedServerVersions         *[]ServerVersionCapabilityV2 `json:"supportedServerVersions,omitempty"`
}
