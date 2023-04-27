package publishedblueprint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceGroupDefinition struct {
	DependsOn *[]string                    `json:"dependsOn,omitempty"`
	Location  *string                      `json:"location,omitempty"`
	Metadata  *ParameterDefinitionMetadata `json:"metadata,omitempty"`
	Name      *string                      `json:"name,omitempty"`
	Tags      *map[string]string           `json:"tags,omitempty"`
}
