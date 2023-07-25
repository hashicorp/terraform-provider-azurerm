package dataset

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobContainerProperties struct {
	ContainerName      string  `json:"containerName"`
	DataSetId          *string `json:"dataSetId,omitempty"`
	ResourceGroup      string  `json:"resourceGroup"`
	StorageAccountName string  `json:"storageAccountName"`
	SubscriptionId     string  `json:"subscriptionId"`
}
