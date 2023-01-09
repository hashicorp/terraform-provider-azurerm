package labplan

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LabPlanUpdateProperties struct {
	AllowedRegions             *[]string              `json:"allowedRegions,omitempty"`
	DefaultAutoShutdownProfile *AutoShutdownProfile   `json:"defaultAutoShutdownProfile,omitempty"`
	DefaultConnectionProfile   *ConnectionProfile     `json:"defaultConnectionProfile,omitempty"`
	DefaultNetworkProfile      *LabPlanNetworkProfile `json:"defaultNetworkProfile,omitempty"`
	LinkedLmsInstance          *string                `json:"linkedLmsInstance,omitempty"`
	SharedGalleryId            *string                `json:"sharedGalleryId,omitempty"`
	SupportInfo                *SupportInfo           `json:"supportInfo,omitempty"`
}
