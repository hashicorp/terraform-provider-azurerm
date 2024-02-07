package azuretrafficcollectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureTrafficCollectorPropertiesFormat struct {
	CollectorPolicies *[]ResourceReference `json:"collectorPolicies,omitempty"`
	ProvisioningState *ProvisioningState   `json:"provisioningState,omitempty"`
	VirtualHub        *ResourceReference   `json:"virtualHub,omitempty"`
}
