package privateclouds

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateCloudProperties struct {
	Circuit                      *Circuit                       `json:"circuit,omitempty"`
	Endpoints                    *Endpoints                     `json:"endpoints,omitempty"`
	IdentitySources              *[]IdentitySource              `json:"identitySources,omitempty"`
	Internet                     *InternetEnum                  `json:"internet,omitempty"`
	ManagementCluster            ManagementCluster              `json:"managementCluster"`
	ManagementNetwork            *string                        `json:"managementNetwork,omitempty"`
	NetworkBlock                 string                         `json:"networkBlock"`
	NsxtCertificateThumbprint    *string                        `json:"nsxtCertificateThumbprint,omitempty"`
	NsxtPassword                 *string                        `json:"nsxtPassword,omitempty"`
	ProvisioningNetwork          *string                        `json:"provisioningNetwork,omitempty"`
	ProvisioningState            *PrivateCloudProvisioningState `json:"provisioningState,omitempty"`
	VcenterCertificateThumbprint *string                        `json:"vcenterCertificateThumbprint,omitempty"`
	VcenterPassword              *string                        `json:"vcenterPassword,omitempty"`
	VmotionNetwork               *string                        `json:"vmotionNetwork,omitempty"`
}
