package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfluentListMetadata struct {
	First     *string `json:"first,omitempty"`
	Last      *string `json:"last,omitempty"`
	Next      *string `json:"next,omitempty"`
	Prev      *string `json:"prev,omitempty"`
	TotalSize *int64  `json:"total_size,omitempty"`
}
