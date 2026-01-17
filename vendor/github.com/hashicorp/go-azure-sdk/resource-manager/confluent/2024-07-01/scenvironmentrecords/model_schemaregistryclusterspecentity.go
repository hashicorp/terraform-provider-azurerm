package scenvironmentrecords

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchemaRegistryClusterSpecEntity struct {
	Cloud        *string                                       `json:"cloud,omitempty"`
	Environment  *SchemaRegistryClusterEnvironmentRegionEntity `json:"environment,omitempty"`
	HTTPEndpoint *string                                       `json:"httpEndpoint,omitempty"`
	Name         *string                                       `json:"name,omitempty"`
	Package      *string                                       `json:"package,omitempty"`
	Region       *SchemaRegistryClusterEnvironmentRegionEntity `json:"region,omitempty"`
}
