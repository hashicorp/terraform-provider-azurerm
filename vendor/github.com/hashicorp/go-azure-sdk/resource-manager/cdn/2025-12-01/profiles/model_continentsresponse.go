package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContinentsResponse struct {
	Continents       *[]ContinentsResponseContinentsItem       `json:"continents,omitempty"`
	CountryOrRegions *[]ContinentsResponseCountryOrRegionsItem `json:"countryOrRegions,omitempty"`
}
