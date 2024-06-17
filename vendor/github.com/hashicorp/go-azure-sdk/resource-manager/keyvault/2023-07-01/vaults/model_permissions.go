package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Permissions struct {
	Certificates *[]CertificatePermissions `json:"certificates,omitempty"`
	Keys         *[]KeyPermissions         `json:"keys,omitempty"`
	Secrets      *[]SecretPermissions      `json:"secrets,omitempty"`
	Storage      *[]StoragePermissions     `json:"storage,omitempty"`
}
