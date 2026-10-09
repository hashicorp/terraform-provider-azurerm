package skillsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchIndexerKnowledgeStoreTableProjectionSelector struct {
	GeneratedKeyName *string                   `json:"generatedKeyName,omitempty"`
	Inputs           *[]InputFieldMappingEntry `json:"inputs,omitempty"`
	ReferenceKeyName *string                   `json:"referenceKeyName,omitempty"`
	Source           *string                   `json:"source,omitempty"`
	SourceContext    *string                   `json:"sourceContext,omitempty"`
	TableName        string                    `json:"tableName"`
}
