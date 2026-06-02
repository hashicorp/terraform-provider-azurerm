package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFileShareConfiguration struct {
	AccountKey        string  `json:"accountKey"`
	AccountName       string  `json:"accountName"`
	AzureFileURL      string  `json:"azureFileUrl"`
	MountOptions      *string `json:"mountOptions,omitempty"`
	RelativeMountPath string  `json:"relativeMountPath"`
}
