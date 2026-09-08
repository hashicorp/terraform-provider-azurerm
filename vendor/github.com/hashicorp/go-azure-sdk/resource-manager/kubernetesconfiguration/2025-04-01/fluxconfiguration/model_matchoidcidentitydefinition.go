package fluxconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MatchOidcIdentityDefinition struct {
	Issuer  *string `json:"issuer,omitempty"`
	Subject *string `json:"subject,omitempty"`
}
