package afddomains

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AFDDomainProperties struct {
	AzureDnsZone                       *ResourceReference          `json:"azureDnsZone,omitempty"`
	DeploymentStatus                   *DeploymentStatus           `json:"deploymentStatus,omitempty"`
	DomainValidationState              *DomainValidationState      `json:"domainValidationState,omitempty"`
	ExtendedProperties                 *map[string]string          `json:"extendedProperties,omitempty"`
	HostName                           string                      `json:"hostName"`
	PreValidatedCustomDomainResourceId *ResourceReference          `json:"preValidatedCustomDomainResourceId,omitempty"`
	ProfileName                        *string                     `json:"profileName,omitempty"`
	ProvisioningState                  *AfdProvisioningState       `json:"provisioningState,omitempty"`
	TlsSettings                        *AFDDomainHTTPSParameters   `json:"tlsSettings,omitempty"`
	ValidationProperties               *DomainValidationProperties `json:"validationProperties,omitempty"`
}
