package netappaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountProperties struct {
	ActiveDirectories *[]ActiveDirectory `json:"activeDirectories,omitempty"`
	Encryption        *AccountEncryption `json:"encryption"`
	ProvisioningState *string            `json:"provisioningState,omitempty"`
}
