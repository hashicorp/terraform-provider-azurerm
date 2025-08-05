package apioperation

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ParameterContract struct {
	DefaultValue *string                              `json:"defaultValue,omitempty"`
	Description  *string                              `json:"description,omitempty"`
	Examples     *map[string]ParameterExampleContract `json:"examples,omitempty"`
	Name         string                               `json:"name"`
	Required     *bool                                `json:"required,omitempty"`
	SchemaId     *string                              `json:"schemaId,omitempty"`
	Type         string                               `json:"type"`
	TypeName     *string                              `json:"typeName,omitempty"`
	Values       *[]string                            `json:"values,omitempty"`
}
