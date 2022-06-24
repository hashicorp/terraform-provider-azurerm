package availableservicetiers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableServiceTier struct {
	CapacityReservationLevel *int64       `json:"capacityReservationLevel,omitempty"`
	DefaultRetention         *int64       `json:"defaultRetention,omitempty"`
	Enabled                  *bool        `json:"enabled,omitempty"`
	LastSkuUpdate            *string      `json:"lastSkuUpdate,omitempty"`
	MaximumRetention         *int64       `json:"maximumRetention,omitempty"`
	MinimumRetention         *int64       `json:"minimumRetention,omitempty"`
	ServiceTier              *SkuNameEnum `json:"serviceTier,omitempty"`
}
