package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionDetails struct {
	Cmk                     *CustomerManagedKeyDetails `json:"cmk,omitempty"`
	DoubleEncryptionEnabled *bool                      `json:"doubleEncryptionEnabled,omitempty"`
}
