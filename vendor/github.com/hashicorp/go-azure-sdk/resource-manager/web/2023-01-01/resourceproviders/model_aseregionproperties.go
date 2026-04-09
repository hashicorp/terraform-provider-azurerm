package resourceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AseRegionProperties struct {
	AvailableOS   *[]string `json:"availableOS,omitempty"`
	AvailableSku  *[]string `json:"availableSku,omitempty"`
	DedicatedHost *bool     `json:"dedicatedHost,omitempty"`
	DisplayName   *string   `json:"displayName,omitempty"`
	Standard      *bool     `json:"standard,omitempty"`
	ZoneRedundant *bool     `json:"zoneRedundant,omitempty"`
}
