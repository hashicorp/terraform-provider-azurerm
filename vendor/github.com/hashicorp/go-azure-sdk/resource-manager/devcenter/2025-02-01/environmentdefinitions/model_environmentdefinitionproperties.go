package environmentdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentDefinitionProperties struct {
	Description      *string                           `json:"description,omitempty"`
	Parameters       *[]EnvironmentDefinitionParameter `json:"parameters,omitempty"`
	TemplatePath     *string                           `json:"templatePath,omitempty"`
	ValidationStatus *CatalogResourceValidationStatus  `json:"validationStatus,omitempty"`
}
