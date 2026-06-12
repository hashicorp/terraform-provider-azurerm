package policysetdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyDefinitionReference struct {
	DefinitionVersion           *string                          `json:"definitionVersion,omitempty"`
	EffectiveDefinitionVersion  *string                          `json:"effectiveDefinitionVersion,omitempty"`
	GroupNames                  *[]string                        `json:"groupNames,omitempty"`
	LatestDefinitionVersion     *string                          `json:"latestDefinitionVersion,omitempty"`
	Parameters                  *map[string]ParameterValuesValue `json:"parameters,omitempty"`
	PolicyDefinitionId          string                           `json:"policyDefinitionId"`
	PolicyDefinitionReferenceId *string                          `json:"policyDefinitionReferenceId,omitempty"`
}
