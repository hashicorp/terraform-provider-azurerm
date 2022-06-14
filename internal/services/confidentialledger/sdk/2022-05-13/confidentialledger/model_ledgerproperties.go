package confidentialledger

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LedgerProperties struct {
	AadBasedSecurityPrincipals  *[]AADBasedSecurityPrincipal  `json:"aadBasedSecurityPrincipals,omitempty"`
	CertBasedSecurityPrincipals *[]CertBasedSecurityPrincipal `json:"certBasedSecurityPrincipals,omitempty"`
	IdentityServiceUri          *string                       `json:"identityServiceUri,omitempty"`
	LedgerInternalNamespace     *string                       `json:"ledgerInternalNamespace,omitempty"`
	LedgerName                  *string                       `json:"ledgerName,omitempty"`
	LedgerType                  *LedgerType                   `json:"ledgerType,omitempty"`
	LedgerUri                   *string                       `json:"ledgerUri,omitempty"`
	ProvisioningState           *ProvisioningState            `json:"provisioningState,omitempty"`
}
