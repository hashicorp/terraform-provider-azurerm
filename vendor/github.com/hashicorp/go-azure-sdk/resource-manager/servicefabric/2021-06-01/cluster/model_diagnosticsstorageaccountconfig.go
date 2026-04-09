package cluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticsStorageAccountConfig struct {
	BlobEndpoint             string  `json:"blobEndpoint"`
	ProtectedAccountKeyName  string  `json:"protectedAccountKeyName"`
	ProtectedAccountKeyName2 *string `json:"protectedAccountKeyName2,omitempty"`
	QueueEndpoint            string  `json:"queueEndpoint"`
	StorageAccountName       string  `json:"storageAccountName"`
	TableEndpoint            string  `json:"tableEndpoint"`
}
