package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFilesIdentityBasedAuthentication struct {
	ActiveDirectoryProperties *ActiveDirectoryProperties `json:"activeDirectoryProperties,omitempty"`
	DefaultSharePermission    *DefaultSharePermission    `json:"defaultSharePermission,omitempty"`
	DirectoryServiceOptions   DirectoryServiceOptions    `json:"directoryServiceOptions"`
}
