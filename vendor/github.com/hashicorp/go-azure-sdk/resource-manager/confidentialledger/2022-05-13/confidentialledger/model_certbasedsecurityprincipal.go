package confidentialledger

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertBasedSecurityPrincipal struct {
	Cert           *string         `json:"cert,omitempty"`
	LedgerRoleName *LedgerRoleName `json:"ledgerRoleName,omitempty"`
}
