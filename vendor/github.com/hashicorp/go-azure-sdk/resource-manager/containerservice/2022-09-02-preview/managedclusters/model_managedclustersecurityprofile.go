package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterSecurityProfile struct {
	AzureKeyVaultKms          *AzureKeyVaultKms                              `json:"azureKeyVaultKms"`
	CustomCATrustCertificates *[]string                                      `json:"customCATrustCertificates,omitempty"`
	Defender                  *ManagedClusterSecurityProfileDefender         `json:"defender"`
	ImageCleaner              *ManagedClusterSecurityProfileImageCleaner     `json:"imageCleaner"`
	NodeRestriction           *ManagedClusterSecurityProfileNodeRestriction  `json:"nodeRestriction"`
	WorkloadIdentity          *ManagedClusterSecurityProfileWorkloadIdentity `json:"workloadIdentity"`
}
