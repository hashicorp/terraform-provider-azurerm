package policysetdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ParameterDefinitionsValue struct {
	AllowedValues *[]interface{}                     `json:"allowedValues,omitempty"`
	DefaultValue  *interface{}                       `json:"defaultValue,omitempty"`
	Metadata      *ParameterDefinitionsValueMetadata `json:"metadata,omitempty"`
	Schema        *interface{}                       `json:"schema,omitempty"`
	Type          *ParameterType                     `json:"type,omitempty"`
}
