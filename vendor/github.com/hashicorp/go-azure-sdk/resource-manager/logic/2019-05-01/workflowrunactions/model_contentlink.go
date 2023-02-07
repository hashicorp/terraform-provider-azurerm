package workflowrunactions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentLink struct {
	ContentHash    *ContentHash `json:"contentHash,omitempty"`
	ContentSize    *int64       `json:"contentSize,omitempty"`
	ContentVersion *string      `json:"contentVersion,omitempty"`
	Metadata       *interface{} `json:"metadata,omitempty"`
	Uri            *string      `json:"uri,omitempty"`
}
