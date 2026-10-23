package documents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IndexingResult struct {
	ErrorMessage *string `json:"errorMessage,omitempty"`
	Key          string  `json:"key"`
	Status       bool    `json:"status"`
	StatusCode   int64   `json:"statusCode"`
}
