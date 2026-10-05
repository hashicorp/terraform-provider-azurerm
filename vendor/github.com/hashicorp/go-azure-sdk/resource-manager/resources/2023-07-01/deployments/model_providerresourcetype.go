package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProviderResourceType struct {
	Aliases           *[]Alias                    `json:"aliases,omitempty"`
	ApiProfiles       *[]ApiProfile               `json:"apiProfiles,omitempty"`
	ApiVersions       *[]string                   `json:"apiVersions,omitempty"`
	Capabilities      *string                     `json:"capabilities,omitempty"`
	DefaultApiVersion *string                     `json:"defaultApiVersion,omitempty"`
	LocationMappings  *[]ProviderExtendedLocation `json:"locationMappings,omitempty"`
	Locations         *[]string                   `json:"locations,omitempty"`
	Properties        *map[string]string          `json:"properties,omitempty"`
	ResourceType      *string                     `json:"resourceType,omitempty"`
	ZoneMappings      *[]ZoneMapping              `json:"zoneMappings,omitempty"`
}
