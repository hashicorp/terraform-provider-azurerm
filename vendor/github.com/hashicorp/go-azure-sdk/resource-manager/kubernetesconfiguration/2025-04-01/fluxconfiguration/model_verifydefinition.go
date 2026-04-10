package fluxconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VerifyDefinition struct {
	MatchOidcIdentity  *[]MatchOidcIdentityDefinition `json:"matchOidcIdentity,omitempty"`
	Provider           *string                        `json:"provider,omitempty"`
	VerificationConfig *map[string]string             `json:"verificationConfig,omitempty"`
}
