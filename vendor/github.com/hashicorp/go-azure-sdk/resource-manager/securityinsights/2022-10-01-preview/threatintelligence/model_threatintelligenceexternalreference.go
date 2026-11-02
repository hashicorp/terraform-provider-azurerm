package threatintelligence

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreatIntelligenceExternalReference struct {
	Description *string            `json:"description,omitempty"`
	ExternalId  *string            `json:"externalId,omitempty"`
	Hashes      *map[string]string `json:"hashes,omitempty"`
	SourceName  *string            `json:"sourceName,omitempty"`
	Url         *string            `json:"url,omitempty"`
}
