package dataset

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ADLSGen2FolderProperties struct {
	DataSetId          *string `json:"dataSetId,omitempty"`
	FileSystem         string  `json:"fileSystem"`
	FolderPath         string  `json:"folderPath"`
	ResourceGroup      string  `json:"resourceGroup"`
	StorageAccountName string  `json:"storageAccountName"`
	SubscriptionId     string  `json:"subscriptionId"`
}
