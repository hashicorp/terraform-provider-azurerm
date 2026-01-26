package backuprestores

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SASTokenParameter struct {
	StorageResourceUri string  `json:"storageResourceUri"`
	Token              *string `json:"token,omitempty"`
	UseManagedIdentity *bool   `json:"useManagedIdentity,omitempty"`
}
