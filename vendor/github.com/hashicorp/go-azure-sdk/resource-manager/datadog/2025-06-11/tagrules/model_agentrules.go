package tagrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentRules struct {
	EnableAgentMonitoring *bool           `json:"enableAgentMonitoring,omitempty"`
	FilteringTags         *[]FilteringTag `json:"filteringTags,omitempty"`
}
