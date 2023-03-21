package simpolicy

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataNetworkConfiguration struct {
	AdditionalAllowedSessionTypes       *[]PduSessionType        `json:"additionalAllowedSessionTypes,omitempty"`
	AllocationAndRetentionPriorityLevel *int64                   `json:"allocationAndRetentionPriorityLevel,omitempty"`
	AllowedServices                     []ServiceResourceId      `json:"allowedServices"`
	DataNetwork                         DataNetworkResourceId    `json:"dataNetwork"`
	DefaultSessionType                  *PduSessionType          `json:"defaultSessionType,omitempty"`
	Fiveqi                              *int64                   `json:"5qi,omitempty"`
	MaximumNumberOfBufferedPackets      *int64                   `json:"maximumNumberOfBufferedPackets,omitempty"`
	PreemptionCapability                *PreemptionCapability    `json:"preemptionCapability,omitempty"`
	PreemptionVulnerability             *PreemptionVulnerability `json:"preemptionVulnerability,omitempty"`
	SessionAmbr                         Ambr                     `json:"sessionAmbr"`
}
