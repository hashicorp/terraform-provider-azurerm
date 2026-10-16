package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileSystemHTTPLogsConfig struct {
	Enabled         *bool  `json:"enabled,omitempty"`
	RetentionInDays *int64 `json:"retentionInDays,omitempty"`
	RetentionInMb   *int64 `json:"retentionInMb,omitempty"`
}
