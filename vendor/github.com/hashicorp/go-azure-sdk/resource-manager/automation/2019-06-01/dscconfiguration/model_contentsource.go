package dscconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentSource struct {
	Hash    *ContentHash       `json:"hash,omitempty"`
	Type    *ContentSourceType `json:"type,omitempty"`
	Value   *string            `json:"value,omitempty"`
	Version *string            `json:"version,omitempty"`
}
