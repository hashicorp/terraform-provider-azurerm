package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IncludedPath struct {
	Indexes *[]Indexes `json:"indexes,omitempty"`
	Path    *string    `json:"path,omitempty"`
}
