package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VectorIndex struct {
	Path string          `json:"path"`
	Type VectorIndexType `json:"type"`
}
