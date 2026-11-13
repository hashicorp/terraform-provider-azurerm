package openapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FleetspaceProperties struct {
	DataRegions                 *[]string                                        `json:"dataRegions,omitempty"`
	FleetspaceApiKind           *FleetspacePropertiesFleetspaceApiKind           `json:"fleetspaceApiKind,omitempty"`
	ProvisioningState           *Status                                          `json:"provisioningState,omitempty"`
	ServiceTier                 *FleetspacePropertiesServiceTier                 `json:"serviceTier,omitempty"`
	ThroughputPoolConfiguration *FleetspacePropertiesThroughputPoolConfiguration `json:"throughputPoolConfiguration,omitempty"`
}
