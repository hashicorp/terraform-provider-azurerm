package storage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SasDefinitionUpdateParameters struct {
	Attributes     *SasDefinitionAttributes `json:"attributes,omitempty"`
	SasType        *SasTokenType            `json:"sasType,omitempty"`
	Tags           *map[string]string       `json:"tags,omitempty"`
	TemplateUri    *string                  `json:"templateUri,omitempty"`
	ValidityPeriod *string                  `json:"validityPeriod,omitempty"`
}
