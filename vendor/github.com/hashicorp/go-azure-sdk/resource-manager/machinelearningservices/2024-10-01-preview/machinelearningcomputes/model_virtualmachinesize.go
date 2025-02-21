package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineSize struct {
	EstimatedVMPrices     *EstimatedVMPrices `json:"estimatedVMPrices,omitempty"`
	Family                *string            `json:"family,omitempty"`
	Gpus                  *int64             `json:"gpus,omitempty"`
	LowPriorityCapable    *bool              `json:"lowPriorityCapable,omitempty"`
	MaxResourceVolumeMB   *int64             `json:"maxResourceVolumeMB,omitempty"`
	MemoryGB              *float64           `json:"memoryGB,omitempty"`
	Name                  *string            `json:"name,omitempty"`
	OsVhdSizeMB           *int64             `json:"osVhdSizeMB,omitempty"`
	PremiumIO             *bool              `json:"premiumIO,omitempty"`
	SupportedComputeTypes *[]string          `json:"supportedComputeTypes,omitempty"`
	VCPUs                 *int64             `json:"vCPUs,omitempty"`
}
