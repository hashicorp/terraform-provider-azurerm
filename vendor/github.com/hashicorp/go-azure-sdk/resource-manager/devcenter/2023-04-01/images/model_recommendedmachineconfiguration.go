package images

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecommendedMachineConfiguration struct {
	Memory *ResourceRange `json:"memory,omitempty"`
	VCPUs  *ResourceRange `json:"vCPUs,omitempty"`
}
