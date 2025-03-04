package netappaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountProperties struct {
	ActiveDirectories *[]ActiveDirectory `json:"activeDirectories,omitempty"`
	DisableShowmount  *bool              `json:"disableShowmount,omitempty"`
	Encryption        *AccountEncryption `json:"encryption,omitempty"`
	ProvisioningState *string            `json:"provisioningState,omitempty"`
}
