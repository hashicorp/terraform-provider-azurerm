package certificateobjectglobalrulestack

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateObject struct {
	AuditComment                *string            `json:"auditComment,omitempty"`
	CertificateSelfSigned       BooleanEnum        `json:"certificateSelfSigned"`
	CertificateSignerResourceId *string            `json:"certificateSignerResourceId,omitempty"`
	Description                 *string            `json:"description,omitempty"`
	Etag                        *string            `json:"etag,omitempty"`
	ProvisioningState           *ProvisioningState `json:"provisioningState,omitempty"`
}
