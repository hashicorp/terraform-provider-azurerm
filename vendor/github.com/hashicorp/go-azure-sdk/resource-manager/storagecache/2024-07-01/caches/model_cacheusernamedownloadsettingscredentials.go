package caches

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CacheUsernameDownloadSettingsCredentials struct {
	BindDn       *string `json:"bindDn,omitempty"`
	BindPassword *string `json:"bindPassword,omitempty"`
}
