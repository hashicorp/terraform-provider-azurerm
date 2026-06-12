package backend

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendTlsProperties struct {
	ValidateCertificateChain *bool `json:"validateCertificateChain,omitempty"`
	ValidateCertificateName  *bool `json:"validateCertificateName,omitempty"`
}
