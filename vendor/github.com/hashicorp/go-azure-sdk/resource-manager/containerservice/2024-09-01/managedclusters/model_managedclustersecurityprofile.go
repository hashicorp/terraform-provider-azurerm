package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterSecurityProfile struct {
	AzureKeyVaultKms *AzureKeyVaultKms                              `json:"azureKeyVaultKms,omitempty"`
	Defender         *ManagedClusterSecurityProfileDefender         `json:"defender,omitempty"`
	ImageCleaner     *ManagedClusterSecurityProfileImageCleaner     `json:"imageCleaner,omitempty"`
	WorkloadIdentity *ManagedClusterSecurityProfileWorkloadIdentity `json:"workloadIdentity,omitempty"`
}
