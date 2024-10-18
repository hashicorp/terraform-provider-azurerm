package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMInfo struct {
	AgentStatus  *string `json:"agentStatus,omitempty"`
	AgentVersion *string `json:"agentVersion,omitempty"`
	VMId         *string `json:"vmId,omitempty"`
}
