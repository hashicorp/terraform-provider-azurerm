package managedenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedCertificateProperties struct {
	DomainControlValidation *ManagedCertificateDomainControlValidation `json:"domainControlValidation,omitempty"`
	Error                   *string                                    `json:"error,omitempty"`
	ProvisioningState       *CertificateProvisioningState              `json:"provisioningState,omitempty"`
	SubjectName             *string                                    `json:"subjectName,omitempty"`
	ValidationToken         *string                                    `json:"validationToken,omitempty"`
}
