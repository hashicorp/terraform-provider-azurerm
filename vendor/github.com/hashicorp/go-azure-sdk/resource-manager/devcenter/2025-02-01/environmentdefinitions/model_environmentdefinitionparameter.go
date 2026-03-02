package environmentdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentDefinitionParameter struct {
	Description *string        `json:"description,omitempty"`
	Id          *string        `json:"id,omitempty"`
	Name        *string        `json:"name,omitempty"`
	ReadOnly    *bool          `json:"readOnly,omitempty"`
	Required    *bool          `json:"required,omitempty"`
	Type        *ParameterType `json:"type,omitempty"`
}
