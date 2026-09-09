package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureStorageInfoValue struct {
	AccessKey   *string               `json:"accessKey,omitempty"`
	AccountName *string               `json:"accountName,omitempty"`
	MountPath   *string               `json:"mountPath,omitempty"`
	Protocol    *AzureStorageProtocol `json:"protocol,omitempty"`
	ShareName   *string               `json:"shareName,omitempty"`
	State       *AzureStorageState    `json:"state,omitempty"`
	Type        *AzureStorageType     `json:"type,omitempty"`
}
