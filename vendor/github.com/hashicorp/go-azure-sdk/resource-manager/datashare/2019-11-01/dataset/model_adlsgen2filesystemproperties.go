package dataset

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ADLSGen2FileSystemProperties struct {
	DataSetId          *string `json:"dataSetId,omitempty"`
	FileSystem         string  `json:"fileSystem"`
	ResourceGroup      string  `json:"resourceGroup"`
	StorageAccountName string  `json:"storageAccountName"`
	SubscriptionId     string  `json:"subscriptionId"`
}
