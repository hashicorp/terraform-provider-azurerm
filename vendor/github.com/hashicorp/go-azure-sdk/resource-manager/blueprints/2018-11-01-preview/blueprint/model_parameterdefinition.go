package blueprint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ParameterDefinition struct {
	AllowedValues *[]interface{}               `json:"allowedValues,omitempty"`
	DefaultValue  *interface{}                 `json:"defaultValue,omitempty"`
	Metadata      *ParameterDefinitionMetadata `json:"metadata,omitempty"`
	Type          TemplateParameterType        `json:"type"`
}
