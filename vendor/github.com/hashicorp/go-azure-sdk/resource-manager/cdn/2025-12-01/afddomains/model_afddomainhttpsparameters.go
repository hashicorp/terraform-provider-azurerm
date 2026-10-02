package afddomains

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AFDDomainHTTPSParameters struct {
	CertificateType          AfdCertificateType                      `json:"certificateType"`
	CipherSuiteSetType       *AfdCipherSuiteSetType                  `json:"cipherSuiteSetType,omitempty"`
	CustomizedCipherSuiteSet *AFDDomainHTTPSCustomizedCipherSuiteSet `json:"customizedCipherSuiteSet,omitempty"`
	MinimumTlsVersion        *AfdMinimumTlsVersion                   `json:"minimumTlsVersion,omitempty"`
	Secret                   *ResourceReference                      `json:"secret,omitempty"`
}
