package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryCredentials struct {
	Identity          *string `json:"identity,omitempty"`
	PasswordSecretRef *string `json:"passwordSecretRef,omitempty"`
	Server            *string `json:"server,omitempty"`
	Username          *string `json:"username,omitempty"`
}
