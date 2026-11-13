package agents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionConfiguration struct {
	AccessLevel *AgentAccessLevel `json:"accessLevel,omitempty"`
	Identity    *string           `json:"identity,omitempty"`
	Mode        *AgentMode        `json:"mode,omitempty"`
}
