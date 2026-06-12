package policydefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyDefinitionProperties struct {
	Description *string                               `json:"description,omitempty"`
	DisplayName *string                               `json:"displayName,omitempty"`
	Metadata    *interface{}                          `json:"metadata,omitempty"`
	Mode        *string                               `json:"mode,omitempty"`
	Parameters  *map[string]ParameterDefinitionsValue `json:"parameters,omitempty"`
	PolicyRule  *interface{}                          `json:"policyRule,omitempty"`
	PolicyType  *PolicyType                           `json:"policyType,omitempty"`
	Version     *string                               `json:"version,omitempty"`
	Versions    *[]string                             `json:"versions,omitempty"`
}
