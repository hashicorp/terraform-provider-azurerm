package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDomainConfiguration struct {
	CertificateURL           *string                      `json:"certificateUrl,omitempty"`
	ExpectedTxtRecordName    *string                      `json:"expectedTxtRecordName,omitempty"`
	ExpectedTxtRecordValue   *string                      `json:"expectedTxtRecordValue,omitempty"`
	FullyQualifiedDomainName string                       `json:"fullyQualifiedDomainName"`
	Identity                 *CustomDomainIdentity        `json:"identity,omitempty"`
	ValidationState          *CustomDomainValidationState `json:"validationState,omitempty"`
}
