package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobOutputDataSourceProperties struct {
	AuthenticationMode *AuthenticationMode `json:"authenticationMode,omitempty"`
	BlobPathPrefix     *string             `json:"blobPathPrefix,omitempty"`
	BlobWriteMode      *BlobWriteMode      `json:"blobWriteMode,omitempty"`
	Container          *string             `json:"container,omitempty"`
	DateFormat         *string             `json:"dateFormat,omitempty"`
	PathPattern        *string             `json:"pathPattern,omitempty"`
	StorageAccounts    *[]StorageAccount   `json:"storageAccounts,omitempty"`
	TimeFormat         *string             `json:"timeFormat,omitempty"`
}
