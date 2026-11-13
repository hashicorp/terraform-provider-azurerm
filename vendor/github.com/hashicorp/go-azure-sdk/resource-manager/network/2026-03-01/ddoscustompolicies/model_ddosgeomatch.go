package ddoscustompolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DdosGeoMatch struct {
	Continent   *DdosContinent `json:"continent,omitempty"`
	CountryCode *string        `json:"countryCode,omitempty"`
}
