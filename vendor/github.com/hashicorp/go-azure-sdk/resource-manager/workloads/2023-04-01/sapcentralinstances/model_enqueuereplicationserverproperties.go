package sapcentralinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnqueueReplicationServerProperties struct {
	ErsVersion    *EnqueueReplicationServerType `json:"ersVersion,omitempty"`
	Health        *SAPHealthState               `json:"health,omitempty"`
	Hostname      *string                       `json:"hostname,omitempty"`
	IPAddress     *string                       `json:"ipAddress,omitempty"`
	InstanceNo    *string                       `json:"instanceNo,omitempty"`
	KernelPatch   *string                       `json:"kernelPatch,omitempty"`
	KernelVersion *string                       `json:"kernelVersion,omitempty"`
}
