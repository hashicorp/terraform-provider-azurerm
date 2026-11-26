package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataLakeStorageAccountDetails struct {
	AccountURL                   *string `json:"accountUrl,omitempty"`
	CreateManagedPrivateEndpoint *bool   `json:"createManagedPrivateEndpoint,omitempty"`
	Filesystem                   *string `json:"filesystem,omitempty"`
	ResourceId                   *string `json:"resourceId,omitempty"`
}
