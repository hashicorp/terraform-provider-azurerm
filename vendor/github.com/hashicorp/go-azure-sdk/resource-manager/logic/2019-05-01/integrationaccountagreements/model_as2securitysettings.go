package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AS2SecuritySettings struct {
	EnableNRRForInboundDecodedMessages  bool    `json:"enableNRRForInboundDecodedMessages"`
	EnableNRRForInboundEncodedMessages  bool    `json:"enableNRRForInboundEncodedMessages"`
	EnableNRRForInboundMDN              bool    `json:"enableNRRForInboundMDN"`
	EnableNRRForOutboundDecodedMessages bool    `json:"enableNRRForOutboundDecodedMessages"`
	EnableNRRForOutboundEncodedMessages bool    `json:"enableNRRForOutboundEncodedMessages"`
	EnableNRRForOutboundMDN             bool    `json:"enableNRRForOutboundMDN"`
	EncryptionCertificateName           *string `json:"encryptionCertificateName,omitempty"`
	OverrideGroupSigningCertificate     bool    `json:"overrideGroupSigningCertificate"`
	Sha2AlgorithmFormat                 *string `json:"sha2AlgorithmFormat,omitempty"`
	SigningCertificateName              *string `json:"signingCertificateName,omitempty"`
}
