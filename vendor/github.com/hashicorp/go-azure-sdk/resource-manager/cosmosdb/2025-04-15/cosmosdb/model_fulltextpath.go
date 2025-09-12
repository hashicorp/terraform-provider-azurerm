package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FullTextPath struct {
	Language *string `json:"language,omitempty"`
	Path     string  `json:"path"`
}
