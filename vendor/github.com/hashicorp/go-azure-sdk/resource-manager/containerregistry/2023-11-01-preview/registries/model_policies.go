package registries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Policies struct {
	AzureADAuthenticationAsArmPolicy *AzureADAuthenticationAsArmPolicy `json:"azureADAuthenticationAsArmPolicy,omitempty"`
	ExportPolicy                     *ExportPolicy                     `json:"exportPolicy,omitempty"`
	QuarantinePolicy                 *QuarantinePolicy                 `json:"quarantinePolicy,omitempty"`
	RetentionPolicy                  *RetentionPolicy                  `json:"retentionPolicy,omitempty"`
	SoftDeletePolicy                 *SoftDeletePolicy                 `json:"softDeletePolicy,omitempty"`
	TrustPolicy                      *TrustPolicy                      `json:"trustPolicy,omitempty"`
}
