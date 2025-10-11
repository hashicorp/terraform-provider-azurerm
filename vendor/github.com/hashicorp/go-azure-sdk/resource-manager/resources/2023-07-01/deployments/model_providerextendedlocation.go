package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProviderExtendedLocation struct {
	ExtendedLocations *[]string `json:"extendedLocations,omitempty"`
	Location          *string   `json:"location,omitempty"`
	Type              *string   `json:"type,omitempty"`
}
