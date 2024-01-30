package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobReferenceInputDataSourceProperties struct {
	AuthenticationMode       *AuthenticationMode `json:"authenticationMode,omitempty"`
	BlobName                 *string             `json:"blobName,omitempty"`
	Container                *string             `json:"container,omitempty"`
	DateFormat               *string             `json:"dateFormat,omitempty"`
	DeltaPathPattern         *string             `json:"deltaPathPattern,omitempty"`
	DeltaSnapshotRefreshRate *string             `json:"deltaSnapshotRefreshRate,omitempty"`
	FullSnapshotRefreshRate  *string             `json:"fullSnapshotRefreshRate,omitempty"`
	PathPattern              *string             `json:"pathPattern,omitempty"`
	SourcePartitionCount     *int64              `json:"sourcePartitionCount,omitempty"`
	StorageAccounts          *[]StorageAccount   `json:"storageAccounts,omitempty"`
	TimeFormat               *string             `json:"timeFormat,omitempty"`
}
