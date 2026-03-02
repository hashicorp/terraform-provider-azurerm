package signalr

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDomainProperties struct {
	CustomCertificate ResourceReference  `json:"customCertificate"`
	DomainName        string             `json:"domainName"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
