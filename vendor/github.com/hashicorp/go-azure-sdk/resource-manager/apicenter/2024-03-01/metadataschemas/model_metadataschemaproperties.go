package metadataschemas

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetadataSchemaProperties struct {
	AssignedTo *[]MetadataAssignment `json:"assignedTo,omitempty"`
	Schema     string                `json:"schema"`
}
