package replicationfabrics

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataStore struct {
	Capacity     *string `json:"capacity,omitempty"`
	FreeSpace    *string `json:"freeSpace,omitempty"`
	SymbolicName *string `json:"symbolicName,omitempty"`
	Type         *string `json:"type,omitempty"`
	Uuid         *string `json:"uuid,omitempty"`
}
