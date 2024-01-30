package resourceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BillingMeterProperties struct {
	BillingLocation *string  `json:"billingLocation,omitempty"`
	FriendlyName    *string  `json:"friendlyName,omitempty"`
	MeterId         *string  `json:"meterId,omitempty"`
	Multiplier      *float64 `json:"multiplier,omitempty"`
	OsType          *string  `json:"osType,omitempty"`
	ResourceType    *string  `json:"resourceType,omitempty"`
	ShortName       *string  `json:"shortName,omitempty"`
}
