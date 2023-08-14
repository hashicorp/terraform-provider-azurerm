package topictypes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicTypeProperties struct {
	Description              *string                     `json:"description,omitempty"`
	DisplayName              *string                     `json:"displayName,omitempty"`
	Provider                 *string                     `json:"provider,omitempty"`
	ProvisioningState        *TopicTypeProvisioningState `json:"provisioningState,omitempty"`
	ResourceRegionType       *ResourceRegionType         `json:"resourceRegionType,omitempty"`
	SourceResourceFormat     *string                     `json:"sourceResourceFormat,omitempty"`
	SupportedLocations       *[]string                   `json:"supportedLocations,omitempty"`
	SupportedScopesForSource *[]TopicTypeSourceScope     `json:"supportedScopesForSource,omitempty"`
}
