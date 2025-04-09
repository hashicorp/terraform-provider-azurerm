package managedenvironmentsstorages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFileProperties struct {
	AccessMode  *AccessMode `json:"accessMode,omitempty"`
	AccountKey  *string     `json:"accountKey,omitempty"`
	AccountName *string     `json:"accountName,omitempty"`
	ShareName   *string     `json:"shareName,omitempty"`
}
