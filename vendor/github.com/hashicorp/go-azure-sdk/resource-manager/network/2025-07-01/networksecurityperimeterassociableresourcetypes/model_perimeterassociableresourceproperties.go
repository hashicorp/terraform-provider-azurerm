package networksecurityperimeterassociableresourcetypes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PerimeterAssociableResourceProperties struct {
	Description       *string            `json:"description,omitempty"`
	DisplayName       *string            `json:"displayName,omitempty"`
	OutboundSupported *bool              `json:"outboundSupported,omitempty"`
	PublicDnsZones    *[]string          `json:"publicDnsZones,omitempty"`
	ReadinessState    *NspReadinessState `json:"readinessState,omitempty"`
	ResourceType      *string            `json:"resourceType,omitempty"`
	ServiceTags       *[]string          `json:"serviceTags,omitempty"`
}
