package policysetdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicySetDefinitionProperties struct {
	Description            *string                               `json:"description,omitempty"`
	DisplayName            *string                               `json:"displayName,omitempty"`
	Metadata               *interface{}                          `json:"metadata,omitempty"`
	Parameters             *map[string]ParameterDefinitionsValue `json:"parameters,omitempty"`
	PolicyDefinitionGroups *[]PolicyDefinitionGroup              `json:"policyDefinitionGroups,omitempty"`
	PolicyDefinitions      []PolicyDefinitionReference           `json:"policyDefinitions"`
	PolicyType             *PolicyType                           `json:"policyType,omitempty"`
	Version                *string                               `json:"version,omitempty"`
	Versions               *[]string                             `json:"versions,omitempty"`
}
