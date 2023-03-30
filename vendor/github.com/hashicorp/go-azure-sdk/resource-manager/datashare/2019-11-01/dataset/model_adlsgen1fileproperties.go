package dataset

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ADLSGen1FileProperties struct {
	AccountName    string  `json:"accountName"`
	DataSetId      *string `json:"dataSetId,omitempty"`
	FileName       string  `json:"fileName"`
	FolderPath     string  `json:"folderPath"`
	ResourceGroup  string  `json:"resourceGroup"`
	SubscriptionId string  `json:"subscriptionId"`
}
