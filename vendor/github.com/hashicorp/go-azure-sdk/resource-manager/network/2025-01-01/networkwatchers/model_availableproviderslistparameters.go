package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableProvidersListParameters struct {
	AzureLocations *[]string `json:"azureLocations,omitempty"`
	City           *string   `json:"city,omitempty"`
	Country        *string   `json:"country,omitempty"`
	State          *string   `json:"state,omitempty"`
}
