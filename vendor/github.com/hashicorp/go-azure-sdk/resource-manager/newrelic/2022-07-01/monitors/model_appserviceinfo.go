package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServiceInfo struct {
	AgentStatus     *string `json:"agentStatus,omitempty"`
	AgentVersion    *string `json:"agentVersion,omitempty"`
	AzureResourceId *string `json:"azureResourceId,omitempty"`
}
