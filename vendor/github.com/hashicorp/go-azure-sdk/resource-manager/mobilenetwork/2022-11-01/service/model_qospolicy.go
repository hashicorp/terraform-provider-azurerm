package service

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QosPolicy struct {
	AllocationAndRetentionPriorityLevel *int64                   `json:"allocationAndRetentionPriorityLevel,omitempty"`
	Fiveqi                              *int64                   `json:"5qi,omitempty"`
	MaximumBitRate                      Ambr                     `json:"maximumBitRate"`
	PreemptionCapability                *PreemptionCapability    `json:"preemptionCapability,omitempty"`
	PreemptionVulnerability             *PreemptionVulnerability `json:"preemptionVulnerability,omitempty"`
}
