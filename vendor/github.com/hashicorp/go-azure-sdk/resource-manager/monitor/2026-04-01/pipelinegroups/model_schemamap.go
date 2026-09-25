package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchemaMap struct {
	RecordMap   []RecordMap    `json:"recordMap"`
	ResourceMap *[]ResourceMap `json:"resourceMap,omitempty"`
	ScopeMap    *[]ScopeMap    `json:"scopeMap,omitempty"`
}
