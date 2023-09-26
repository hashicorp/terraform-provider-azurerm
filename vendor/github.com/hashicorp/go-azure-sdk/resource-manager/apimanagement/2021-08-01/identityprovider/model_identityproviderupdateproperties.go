package identityprovider

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IdentityProviderUpdateProperties struct {
	AllowedTenants           *[]string             `json:"allowedTenants,omitempty"`
	Authority                *string               `json:"authority,omitempty"`
	ClientId                 *string               `json:"clientId,omitempty"`
	ClientSecret             *string               `json:"clientSecret,omitempty"`
	PasswordResetPolicyName  *string               `json:"passwordResetPolicyName,omitempty"`
	ProfileEditingPolicyName *string               `json:"profileEditingPolicyName,omitempty"`
	SigninPolicyName         *string               `json:"signinPolicyName,omitempty"`
	SigninTenant             *string               `json:"signinTenant,omitempty"`
	SignupPolicyName         *string               `json:"signupPolicyName,omitempty"`
	Type                     *IdentityProviderType `json:"type,omitempty"`
}
