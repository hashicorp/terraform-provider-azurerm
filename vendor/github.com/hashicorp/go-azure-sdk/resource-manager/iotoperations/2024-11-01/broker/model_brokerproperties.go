package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BrokerProperties struct {
	Advanced                *AdvancedSettings        `json:"advanced,omitempty"`
	Cardinality             *Cardinality             `json:"cardinality,omitempty"`
	Diagnostics             *BrokerDiagnostics       `json:"diagnostics,omitempty"`
	DiskBackedMessageBuffer *DiskBackedMessageBuffer `json:"diskBackedMessageBuffer,omitempty"`
	GenerateResourceLimits  *GenerateResourceLimits  `json:"generateResourceLimits,omitempty"`
	MemoryProfile           *BrokerMemoryProfile     `json:"memoryProfile,omitempty"`
	ProvisioningState       *ProvisioningState       `json:"provisioningState,omitempty"`
}
