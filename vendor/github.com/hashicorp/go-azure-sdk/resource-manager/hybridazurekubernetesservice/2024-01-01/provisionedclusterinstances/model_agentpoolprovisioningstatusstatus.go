package provisionedclusterinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentPoolProvisioningStatusStatus struct {
	CurrentState  *ResourceProvisioningState `json:"currentState,omitempty"`
	ErrorMessage  *string                    `json:"errorMessage,omitempty"`
	ReadyReplicas *[]AgentPoolUpdateProfile  `json:"readyReplicas,omitempty"`
}
