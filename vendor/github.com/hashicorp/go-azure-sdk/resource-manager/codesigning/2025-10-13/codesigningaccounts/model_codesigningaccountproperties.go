package codesigningaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CodeSigningAccountProperties struct {
	AccountUri        *string            `json:"accountUri,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Sku               *AccountSku        `json:"sku,omitempty"`
}
