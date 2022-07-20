package eventhubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DestinationProperties struct {
	ArchiveNameFormat        *string `json:"archiveNameFormat,omitempty"`
	BlobContainer            *string `json:"blobContainer,omitempty"`
	StorageAccountResourceId *string `json:"storageAccountResourceId,omitempty"`
}
