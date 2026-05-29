package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceUploadDefinition struct {
	RelativePath *string `json:"relativePath,omitempty"`
	UploadURL    *string `json:"uploadUrl,omitempty"`
}
