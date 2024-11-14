package metadataschemas

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetadataAssignment struct {
	Deprecated *bool                     `json:"deprecated,omitempty"`
	Entity     *MetadataAssignmentEntity `json:"entity,omitempty"`
	Required   *bool                     `json:"required,omitempty"`
}
