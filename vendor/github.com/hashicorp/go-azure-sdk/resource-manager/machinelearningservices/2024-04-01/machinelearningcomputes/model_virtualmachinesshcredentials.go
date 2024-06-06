package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineSshCredentials struct {
	Password       *string `json:"password,omitempty"`
	PrivateKeyData *string `json:"privateKeyData,omitempty"`
	PublicKeyData  *string `json:"publicKeyData,omitempty"`
	Username       *string `json:"username,omitempty"`
}
