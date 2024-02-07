package templatespecversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TemplateSpecVersionProperties struct {
	Description      *string                   `json:"description,omitempty"`
	LinkedTemplates  *[]LinkedTemplateArtifact `json:"linkedTemplates,omitempty"`
	MainTemplate     *interface{}              `json:"mainTemplate,omitempty"`
	Metadata         *interface{}              `json:"metadata,omitempty"`
	UiFormDefinition *interface{}              `json:"uiFormDefinition,omitempty"`
}
