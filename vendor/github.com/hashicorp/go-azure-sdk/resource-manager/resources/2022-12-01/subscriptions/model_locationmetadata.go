package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationMetadata struct {
	Geography        *string         `json:"geography,omitempty"`
	GeographyGroup   *string         `json:"geographyGroup,omitempty"`
	HomeLocation     *string         `json:"homeLocation,omitempty"`
	Latitude         *string         `json:"latitude,omitempty"`
	Longitude        *string         `json:"longitude,omitempty"`
	PairedRegion     *[]PairedRegion `json:"pairedRegion,omitempty"`
	PhysicalLocation *string         `json:"physicalLocation,omitempty"`
	RegionCategory   *RegionCategory `json:"regionCategory,omitempty"`
	RegionType       *RegionType     `json:"regionType,omitempty"`
}
