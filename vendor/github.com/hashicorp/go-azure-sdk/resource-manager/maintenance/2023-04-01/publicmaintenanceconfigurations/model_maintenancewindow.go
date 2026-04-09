package publicmaintenanceconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceWindow struct {
	Duration           *string `json:"duration,omitempty"`
	ExpirationDateTime *string `json:"expirationDateTime,omitempty"`
	RecurEvery         *string `json:"recurEvery,omitempty"`
	StartDateTime      *string `json:"startDateTime,omitempty"`
	TimeZone           *string `json:"timeZone,omitempty"`
}
