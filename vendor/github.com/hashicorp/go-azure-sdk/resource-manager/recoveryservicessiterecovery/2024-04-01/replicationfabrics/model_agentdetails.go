package replicationfabrics

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentDetails struct {
	AgentId   *string             `json:"agentId,omitempty"`
	BiosId    *string             `json:"biosId,omitempty"`
	Disks     *[]AgentDiskDetails `json:"disks,omitempty"`
	Fqdn      *string             `json:"fqdn,omitempty"`
	MachineId *string             `json:"machineId,omitempty"`
}
