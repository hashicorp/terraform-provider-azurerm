package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificatePolicy struct {
	Attributes      *CertificateAttributes     `json:"attributes,omitempty"`
	Id              *string                    `json:"id,omitempty"`
	Issuer          *IssuerParameters          `json:"issuer,omitempty"`
	KeyProps        *KeyProperties             `json:"key_props,omitempty"`
	LifetimeActions *[]LifetimeAction          `json:"lifetime_actions,omitempty"`
	SecretProps     *SecretProperties          `json:"secret_props,omitempty"`
	X509Props       *X509CertificateProperties `json:"x509_props,omitempty"`
}
