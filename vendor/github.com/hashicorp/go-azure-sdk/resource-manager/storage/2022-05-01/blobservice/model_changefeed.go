package blobservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ChangeFeed struct {
	Enabled         *bool  `json:"enabled,omitempty"`
	RetentionInDays *int64 `json:"retentionInDays,omitempty"`
}
