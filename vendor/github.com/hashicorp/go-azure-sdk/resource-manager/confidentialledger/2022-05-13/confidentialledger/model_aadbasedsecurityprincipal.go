package confidentialledger

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AADBasedSecurityPrincipal struct {
	LedgerRoleName *LedgerRoleName `json:"ledgerRoleName,omitempty"`
	PrincipalId    *string         `json:"principalId,omitempty"`
	TenantId       *string         `json:"tenantId,omitempty"`
}
