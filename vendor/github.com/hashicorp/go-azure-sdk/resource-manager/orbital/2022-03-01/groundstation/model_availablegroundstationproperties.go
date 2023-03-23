package groundstation

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableGroundStationProperties struct {
	AltitudeMeters   *float64     `json:"altitudeMeters,omitempty"`
	City             *string      `json:"city,omitempty"`
	LatitudeDegrees  *float64     `json:"latitudeDegrees,omitempty"`
	LongitudeDegrees *float64     `json:"longitudeDegrees,omitempty"`
	ProviderName     *string      `json:"providerName,omitempty"`
	ReleaseMode      *ReleaseMode `json:"releaseMode,omitempty"`
}
