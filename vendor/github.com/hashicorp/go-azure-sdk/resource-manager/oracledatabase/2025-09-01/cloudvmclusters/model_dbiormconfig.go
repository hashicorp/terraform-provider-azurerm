package cloudvmclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DbIormConfig struct {
	DbName          *string `json:"dbName,omitempty"`
	FlashCacheLimit *string `json:"flashCacheLimit,omitempty"`
	Share           *int64  `json:"share,omitempty"`
}
