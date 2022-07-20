package configurationstores

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationStorePropertiesUpdateParameters struct {
	DisableLocalAuth      *bool                 `json:"disableLocalAuth,omitempty"`
	EnablePurgeProtection *bool                 `json:"enablePurgeProtection,omitempty"`
	Encryption            *EncryptionProperties `json:"encryption,omitempty"`
	PublicNetworkAccess   *PublicNetworkAccess  `json:"publicNetworkAccess,omitempty"`
}
