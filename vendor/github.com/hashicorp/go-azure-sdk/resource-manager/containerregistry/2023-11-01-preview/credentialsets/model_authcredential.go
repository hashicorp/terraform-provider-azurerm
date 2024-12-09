package credentialsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthCredential struct {
	CredentialHealth         *CredentialHealth `json:"credentialHealth,omitempty"`
	Name                     *CredentialName   `json:"name,omitempty"`
	PasswordSecretIdentifier *string           `json:"passwordSecretIdentifier,omitempty"`
	UsernameSecretIdentifier *string           `json:"usernameSecretIdentifier,omitempty"`
}
