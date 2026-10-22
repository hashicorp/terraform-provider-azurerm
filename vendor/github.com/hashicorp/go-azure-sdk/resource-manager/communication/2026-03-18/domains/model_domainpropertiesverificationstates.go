package domains

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DomainPropertiesVerificationStates struct {
	DKIM   *VerificationStatusRecord `json:"DKIM,omitempty"`
	DKIM2  *VerificationStatusRecord `json:"DKIM2,omitempty"`
	DMARC  *VerificationStatusRecord `json:"DMARC,omitempty"`
	Domain *VerificationStatusRecord `json:"Domain,omitempty"`
	SPF    *VerificationStatusRecord `json:"SPF,omitempty"`
}
