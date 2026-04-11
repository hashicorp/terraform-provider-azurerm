package indexers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FieldMapping struct {
	MappingFunction *FieldMappingFunction `json:"mappingFunction,omitempty"`
	SourceFieldName string                `json:"sourceFieldName"`
	TargetFieldName *string               `json:"targetFieldName,omitempty"`
}
