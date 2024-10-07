package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterSecurityProfile struct {
	AzureKeyVaultKms          *AzureKeyVaultKms                              `json:"azureKeyVaultKms,omitempty"`
	CustomCATrustCertificates *[]string                                      `json:"customCATrustCertificates,omitempty"`
	Defender                  *ManagedClusterSecurityProfileDefender         `json:"defender,omitempty"`
	ImageCleaner              *ManagedClusterSecurityProfileImageCleaner     `json:"imageCleaner,omitempty"`
	NodeRestriction           *ManagedClusterSecurityProfileNodeRestriction  `json:"nodeRestriction,omitempty"`
	WorkloadIdentity          *ManagedClusterSecurityProfileWorkloadIdentity `json:"workloadIdentity,omitempty"`
}
