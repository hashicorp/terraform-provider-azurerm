package partnertopics

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InlineEventProperties struct {
	DataSchemaUrl    *string `json:"dataSchemaUrl,omitempty"`
	Description      *string `json:"description,omitempty"`
	DisplayName      *string `json:"displayName,omitempty"`
	DocumentationUrl *string `json:"documentationUrl,omitempty"`
}
