package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SlotConfigNames struct {
	AppSettingNames         *[]string `json:"appSettingNames,omitempty"`
	AzureStorageConfigNames *[]string `json:"azureStorageConfigNames,omitempty"`
	ConnectionStringNames   *[]string `json:"connectionStringNames,omitempty"`
}
