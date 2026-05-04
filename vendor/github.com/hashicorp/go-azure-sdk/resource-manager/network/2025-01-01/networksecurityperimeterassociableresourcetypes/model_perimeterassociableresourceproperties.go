package networksecurityperimeterassociableresourcetypes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PerimeterAssociableResourceProperties struct {
	DisplayName    *string   `json:"displayName,omitempty"`
	PublicDnsZones *[]string `json:"publicDnsZones,omitempty"`
	ResourceType   *string   `json:"resourceType,omitempty"`
}
