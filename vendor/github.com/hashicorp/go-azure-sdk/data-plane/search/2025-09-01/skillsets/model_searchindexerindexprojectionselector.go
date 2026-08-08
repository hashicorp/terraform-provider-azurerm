package skillsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchIndexerIndexProjectionSelector struct {
	Mappings           []InputFieldMappingEntry `json:"mappings"`
	ParentKeyFieldName string                   `json:"parentKeyFieldName"`
	SourceContext      string                   `json:"sourceContext"`
	TargetIndexName    string                   `json:"targetIndexName"`
}
