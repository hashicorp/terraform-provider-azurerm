package domains

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DomainProperties struct {
	DataLocation           *string                              `json:"dataLocation,omitempty"`
	DomainManagement       DomainManagement                     `json:"domainManagement"`
	FromSenderDomain       *string                              `json:"fromSenderDomain,omitempty"`
	MailFromSenderDomain   *string                              `json:"mailFromSenderDomain,omitempty"`
	ProvisioningState      *DomainsProvisioningState            `json:"provisioningState,omitempty"`
	UserEngagementTracking *UserEngagementTracking              `json:"userEngagementTracking,omitempty"`
	VerificationRecords    *DomainPropertiesVerificationRecords `json:"verificationRecords,omitempty"`
	VerificationStates     *DomainPropertiesVerificationStates  `json:"verificationStates,omitempty"`
}
