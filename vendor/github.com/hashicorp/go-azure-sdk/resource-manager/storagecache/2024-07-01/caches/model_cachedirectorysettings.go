package caches

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CacheDirectorySettings struct {
	ActiveDirectory  *CacheActiveDirectorySettings  `json:"activeDirectory,omitempty"`
	UsernameDownload *CacheUsernameDownloadSettings `json:"usernameDownload,omitempty"`
}
