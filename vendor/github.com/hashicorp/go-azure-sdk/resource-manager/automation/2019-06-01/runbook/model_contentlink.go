package runbook

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentLink struct {
	ContentHash *ContentHash `json:"contentHash"`
	Uri         *string      `json:"uri,omitempty"`
	Version     *string      `json:"version,omitempty"`
}
