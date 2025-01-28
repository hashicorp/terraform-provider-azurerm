package exports

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportDeliveryDestination struct {
	Container      string  `json:"container"`
	ResourceId     *string `json:"resourceId,omitempty"`
	RootFolderPath *string `json:"rootFolderPath,omitempty"`
	SasToken       *string `json:"sasToken,omitempty"`
	StorageAccount *string `json:"storageAccount,omitempty"`
}
