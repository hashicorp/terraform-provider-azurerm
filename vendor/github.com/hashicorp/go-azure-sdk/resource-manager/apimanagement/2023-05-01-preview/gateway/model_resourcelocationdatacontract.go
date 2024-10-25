package gateway

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceLocationDataContract struct {
	City            *string `json:"city,omitempty"`
	CountryOrRegion *string `json:"countryOrRegion,omitempty"`
	District        *string `json:"district,omitempty"`
	Name            string  `json:"name"`
}
