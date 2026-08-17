package afdorigins

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AFDOriginProperties struct {
	AzureOrigin                 *ResourceReference                   `json:"azureOrigin,omitempty"`
	DeploymentStatus            *DeploymentStatus                    `json:"deploymentStatus,omitempty"`
	EnabledState                *EnabledState                        `json:"enabledState,omitempty"`
	EnforceCertificateNameCheck *bool                                `json:"enforceCertificateNameCheck,omitempty"`
	HTTPPort                    *int64                               `json:"httpPort,omitempty"`
	HTTPSPort                   *int64                               `json:"httpsPort,omitempty"`
	HostName                    *string                              `json:"hostName,omitempty"`
	OriginGroupName             *string                              `json:"originGroupName,omitempty"`
	OriginHostHeader            *string                              `json:"originHostHeader,omitempty"`
	Priority                    *int64                               `json:"priority,omitempty"`
	ProvisioningState           *AfdProvisioningState                `json:"provisioningState,omitempty"`
	SharedPrivateLinkResource   *SharedPrivateLinkResourceProperties `json:"sharedPrivateLinkResource,omitempty"`
	Weight                      *int64                               `json:"weight,omitempty"`
}
