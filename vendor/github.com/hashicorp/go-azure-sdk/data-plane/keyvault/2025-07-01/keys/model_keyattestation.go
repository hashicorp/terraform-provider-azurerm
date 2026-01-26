package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyAttestation struct {
	CertificatePemFile    *string `json:"certificatePemFile,omitempty"`
	PrivateKeyAttestation *string `json:"privateKeyAttestation,omitempty"`
	PublicKeyAttestation  *string `json:"publicKeyAttestation,omitempty"`
	Version               *string `json:"version,omitempty"`
}
