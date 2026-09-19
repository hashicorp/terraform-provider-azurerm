package skillsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchIndexerKnowledgeStoreProjection struct {
	Files   *[]SearchIndexerKnowledgeStoreBlobProjectionSelector  `json:"files,omitempty"`
	Objects *[]SearchIndexerKnowledgeStoreBlobProjectionSelector  `json:"objects,omitempty"`
	Tables  *[]SearchIndexerKnowledgeStoreTableProjectionSelector `json:"tables,omitempty"`
}
