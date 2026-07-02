package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchMetadataSchema struct {
	Name    *string `json:"name,omitempty"`
	Version *int64  `json:"version,omitempty"`
}
