package connectedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedClusterPatchProperties struct {
	AzureHybridBenefit  *AzureHybridBenefit `json:"azureHybridBenefit,omitempty"`
	Distribution        *string             `json:"distribution,omitempty"`
	DistributionVersion *string             `json:"distributionVersion,omitempty"`
}
