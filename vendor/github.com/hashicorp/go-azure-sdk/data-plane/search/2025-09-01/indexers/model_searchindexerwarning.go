package indexers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchIndexerWarning struct {
	Details           *string `json:"details,omitempty"`
	DocumentationLink *string `json:"documentationLink,omitempty"`
	Key               *string `json:"key,omitempty"`
	Message           string  `json:"message"`
	Name              *string `json:"name,omitempty"`
}
