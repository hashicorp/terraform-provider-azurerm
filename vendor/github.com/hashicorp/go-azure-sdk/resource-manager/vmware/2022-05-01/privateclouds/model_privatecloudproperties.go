package privateclouds

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateCloudProperties struct {
	Availability                 *AvailabilityProperties        `json:"availability,omitempty"`
	Circuit                      *Circuit                       `json:"circuit,omitempty"`
	Encryption                   *Encryption                    `json:"encryption,omitempty"`
	Endpoints                    *Endpoints                     `json:"endpoints,omitempty"`
	ExternalCloudLinks           *[]string                      `json:"externalCloudLinks,omitempty"`
	IdentitySources              *[]IdentitySource              `json:"identitySources,omitempty"`
	Internet                     *InternetEnum                  `json:"internet,omitempty"`
	ManagementCluster            CommonClusterProperties        `json:"managementCluster"`
	ManagementNetwork            *string                        `json:"managementNetwork,omitempty"`
	NetworkBlock                 string                         `json:"networkBlock"`
	NsxPublicIPQuotaRaised       *NsxPublicIPQuotaRaisedEnum    `json:"nsxPublicIpQuotaRaised,omitempty"`
	NsxtCertificateThumbprint    *string                        `json:"nsxtCertificateThumbprint,omitempty"`
	NsxtPassword                 *string                        `json:"nsxtPassword,omitempty"`
	ProvisioningNetwork          *string                        `json:"provisioningNetwork,omitempty"`
	ProvisioningState            *PrivateCloudProvisioningState `json:"provisioningState,omitempty"`
	SecondaryCircuit             *Circuit                       `json:"secondaryCircuit,omitempty"`
	VMotionNetwork               *string                        `json:"vmotionNetwork,omitempty"`
	VcenterCertificateThumbprint *string                        `json:"vcenterCertificateThumbprint,omitempty"`
	VcenterPassword              *string                        `json:"vcenterPassword,omitempty"`
}
