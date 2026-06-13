package indexers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchIndexerError struct {
	Details           *string `json:"details,omitempty"`
	DocumentationLink *string `json:"documentationLink,omitempty"`
	ErrorMessage      string  `json:"errorMessage"`
	Key               *string `json:"key,omitempty"`
	Name              *string `json:"name,omitempty"`
	StatusCode        int64   `json:"statusCode"`
}
