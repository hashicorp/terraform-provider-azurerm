package registries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomRegistryCredentials struct {
	Identity *string       `json:"identity,omitempty"`
	Password *SecretObject `json:"password,omitempty"`
	UserName *SecretObject `json:"userName,omitempty"`
}
