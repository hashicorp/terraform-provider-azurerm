package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableProvidersListCountry struct {
	CountryName *string                        `json:"countryName,omitempty"`
	Providers   *[]string                      `json:"providers,omitempty"`
	States      *[]AvailableProvidersListState `json:"states,omitempty"`
}
