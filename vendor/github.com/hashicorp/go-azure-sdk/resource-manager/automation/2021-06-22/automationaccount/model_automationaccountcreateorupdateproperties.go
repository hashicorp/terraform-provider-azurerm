package automationaccount

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationAccountCreateOrUpdateProperties struct {
	DisableLocalAuth    *bool                 `json:"disableLocalAuth,omitempty"`
	Encryption          *EncryptionProperties `json:"encryption,omitempty"`
	PublicNetworkAccess *bool                 `json:"publicNetworkAccess,omitempty"`
	Sku                 *Sku                  `json:"sku,omitempty"`
}
