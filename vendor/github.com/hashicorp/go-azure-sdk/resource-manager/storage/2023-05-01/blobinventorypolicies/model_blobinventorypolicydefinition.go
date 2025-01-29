package blobinventorypolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobInventoryPolicyDefinition struct {
	Filters      *BlobInventoryPolicyFilter `json:"filters,omitempty"`
	Format       Format                     `json:"format"`
	ObjectType   ObjectType                 `json:"objectType"`
	Schedule     Schedule                   `json:"schedule"`
	SchemaFields []string                   `json:"schemaFields"`
}
