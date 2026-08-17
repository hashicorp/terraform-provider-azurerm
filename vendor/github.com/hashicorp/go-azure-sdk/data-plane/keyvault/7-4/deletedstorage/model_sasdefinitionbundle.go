package deletedstorage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SasDefinitionBundle struct {
	Attributes     *SasDefinitionAttributes `json:"attributes,omitempty"`
	Id             *string                  `json:"id,omitempty"`
	SasType        *SasTokenType            `json:"sasType,omitempty"`
	Sid            *string                  `json:"sid,omitempty"`
	Tags           *map[string]string       `json:"tags,omitempty"`
	TemplateUri    *string                  `json:"templateUri,omitempty"`
	ValidityPeriod *string                  `json:"validityPeriod,omitempty"`
}
