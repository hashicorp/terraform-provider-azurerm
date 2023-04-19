package registries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceUploadDefinition struct {
	RelativePath *string `json:"relativePath,omitempty"`
	UploadUrl    *string `json:"uploadUrl,omitempty"`
}
