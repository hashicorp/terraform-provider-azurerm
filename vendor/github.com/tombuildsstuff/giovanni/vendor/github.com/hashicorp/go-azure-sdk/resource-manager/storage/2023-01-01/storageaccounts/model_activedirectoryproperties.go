package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActiveDirectoryProperties struct {
	AccountType       *AccountType `json:"accountType,omitempty"`
	AzureStorageSid   *string      `json:"azureStorageSid,omitempty"`
	DomainGuid        string       `json:"domainGuid"`
	DomainName        string       `json:"domainName"`
	DomainSid         *string      `json:"domainSid,omitempty"`
	ForestName        *string      `json:"forestName,omitempty"`
	NetBiosDomainName *string      `json:"netBiosDomainName,omitempty"`
	SamAccountName    *string      `json:"samAccountName,omitempty"`
}
