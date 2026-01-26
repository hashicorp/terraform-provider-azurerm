package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateOperation struct {
	CancellationRequested *bool             `json:"cancellation_requested,omitempty"`
	Csr                   *string           `json:"csr,omitempty"`
	Error                 *Error            `json:"error,omitempty"`
	Id                    *string           `json:"id,omitempty"`
	Issuer                *IssuerParameters `json:"issuer,omitempty"`
	PreserveCertOrder     *bool             `json:"preserveCertOrder,omitempty"`
	RequestId             *string           `json:"request_id,omitempty"`
	Status                *string           `json:"status,omitempty"`
	StatusDetails         *string           `json:"status_details,omitempty"`
	Target                *string           `json:"target,omitempty"`
}
