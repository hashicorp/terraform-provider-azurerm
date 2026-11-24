package passwordauthentications

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PasswordAuthenticationPropertiesForPatchUpdate struct {
	AdministratorLoginPassword *string              `json:"administratorLoginPassword,omitempty"`
	State                      *AuthenticationState `json:"state,omitempty"`
}
