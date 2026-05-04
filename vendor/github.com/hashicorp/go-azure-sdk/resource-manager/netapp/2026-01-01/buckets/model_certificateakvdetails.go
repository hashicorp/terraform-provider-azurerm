package buckets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateAkvDetails struct {
	CertificateKeyVaultUri *string `json:"certificateKeyVaultUri,omitempty"`
	CertificateName        *string `json:"certificateName,omitempty"`
}
